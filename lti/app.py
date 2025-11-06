import os, json
import secrets
from flask import Flask, request, jsonify, redirect, session, make_response
from dotenv import load_dotenv
from pylti1p3.contrib.flask import (
    FlaskOIDCLogin, FlaskMessageLaunch, FlaskRequest, FlaskCacheDataStorage
)
from pylti1p3.tool_config import ToolConfDict
from flask_caching import Cache
from werkzeug.middleware.proxy_fix import ProxyFix


load_dotenv(".env")

# ---- Tool Configuration ----- #
TENANT_ISS = os.environ.get("CANVAS_ISSUER", "https://ufl.instructure.com")
GLOBAL_ISS = "https://canvas.instructure.com"  # Canvas sometimes uses this for iss

def require_env(name):
    value = os.getenv(name)
    if not value:
        raise RuntimeError(f"Missing required environment variable: {name}")
    return value

DEPLOY_IDS = require_env("CANVAS_DEPLOYMENT_ID")

DEPLOY_ID = require_env("CANVAS_DEPLOYMENT_ID")

REG = {
    "default": True,
    "client_id": require_env("CANVAS_CLIENT_ID"),
    "auth_login_url": require_env("CANVAS_OIDC_AUTH_URL"),
    "auth_token_url": require_env("CANVAS_TOKEN_URL"),
    "auth_audience": os.getenv("CANVAS_TOKEN_AUDIENCE"),
    "key_set_url": require_env("CANVAS_JWKS_URL"),
    "deployment_ids": [DEPLOY_ID],
}

tool_conf = ToolConfDict({
    TENANT_ISS: [REG],
    GLOBAL_ISS: [REG],   
})

print("Configured issuers:", [TENANT_ISS, GLOBAL_ISS])
print("Configured client_id:", REG["client_id"])
print("Configured deployment_ids:", DEPLOY_IDS)


# ---- Flask App ---- #
app = Flask(__name__)
app.wsgi_app = ProxyFix(app.wsgi_app, x_for=1, x_proto=1, x_host=1, x_port=1)
app.config.update({
    "PREFERRED_URL_SCHEME": "https",
    "SESSION_COOKIE_DOMAIN": None,
    "SESSION_COOKIE_SAMESITE": "None",
    "SESSION_COOKIE_SECURE": True,
    "X_FRAME_OPTIONS": None,           # disables Flask default header
    "TALISMAN_FRAME_OPTIONS": None,    # disables Flask-Talisman frame header
})
app.secret_key = os.environ.get("FLASK_SECRET", "change-me")
TOOL_JWKS = {"keys": []}  # replace with  actual tool JWKS later when we start signing responses


# -- Cache Instance -- #
# create a cache instance (lives in memory for dev gets wiped when Flask restarts) -> for prod we need something that survives across processes like Redis/Memcached
app.config.update({
    "CACHE_TYPE": "RedisCache",
    "CACHE_REDIS_URL": os.environ.get("REDIS_URL", "redis://redis:6379/0"),
    "CACHE_DEFAULT_TIMEOUT": 900  
})
app_cache = Cache(app)
launch_store = FlaskCacheDataStorage(app_cache)

# -- ENDPOINTS -- #
@app.route("/lti/jwks.json")
def jwks():
    # Serve public keys so Canvas can validate any JWTs (e.g., for Deep Linking responses).
    return jsonify(TOOL_JWKS)



def _harden_state_cookie(resp):
    cookies = resp.headers.getlist("Set-Cookie")
    resp.headers.set("Set-Cookie", "")
    for c in cookies:
        if c.startswith("lti1p3-state-"):
            if "SameSite=None" not in c: c += "; SameSite=None"
            if "Secure" not in c: c += "; Secure"
            if "Path=" not in c: c += "; Path=/"
        resp.headers.add("Set-Cookie", c)
    return resp

@app.route("/lti/oidc-login", methods=["GET","POST"])
def oidc_login():
    req = FlaskRequest()
    login = FlaskOIDCLogin(req, tool_conf, launch_data_storage=launch_store)
    target_link_uri = req.get_param("target_link_uri")
    if not target_link_uri: return jsonify({"error":"missing target_link_uri"}), 400
    resp = login.redirect(target_link_uri)
    return _harden_state_cookie(resp)



COOKIE_NAME = "__Host-gw_sid"

def _set_partitioned(resp, name, value, max_age):
    resp.set_cookie(name, value, max_age=max_age, path="/",
                    secure=True, httponly=True, samesite="None")
    # append Partitioned
    cookies = resp.headers.getlist("Set-Cookie")
    resp.headers.set("Set-Cookie", "")
    for c in cookies:
        if c.startswith(f"{name}=") and "Partitioned" not in c:
            c += "; Partitioned"
        resp.headers.add("Set-Cookie", c)
    return resp

@app.route("/lti/launch", methods=["POST"])
def launch():
    req = FlaskRequest()
    ml = FlaskMessageLaunch(req, tool_conf, launch_data_storage=launch_store).validate()

    # issue your own sid and cache launch data keyed by it
    sid = secrets.token_urlsafe(32)
    app_cache.set(f"sid:{sid}", ml.get_launch_data(), timeout=3600)

    resp = make_response(redirect("https://dev.gradewise.org/app"))
    _set_partitioned(resp, COOKIE_NAME, sid, 3600)
    return resp

@app.route("/lti/cookie-init")
def cookie_init():
    ret = request.args.get("return", "https://dev.gradewise.org/app")
    sid = secrets.token_urlsafe(32)
    app_cache.set(f"sid:{sid}", {"boot": True}, timeout=900)
    resp = make_response(redirect(ret))
    _set_partitioned(resp, "__Host-gw_sid", sid, 900)
    return resp

# Just a simple Flask health route to confirm everything is working - not needed for LTI but it's handy to have
@app.route("/lti", methods=["GET","HEAD"])
def health():
    return jsonify({
        "name": "Gradewise LTI",
        "status": "ok",
        "endpoints": ["/jwks.json", "POST /oidc-login", "POST /launch"]
    })

@app.before_request
def _log(): print("LTI hit", request.path)

@app.after_request
def add_csp(resp):
    resp.headers.pop("X-Frame-Options", None)
    resp.headers["Content-Security-Policy"] = (
        "default-src 'self'; "
        "frame-ancestors https://*.instructure.com https://*.canvaslms.com; "
        "script-src 'self' 'unsafe-inline'; "
        "style-src 'self' 'unsafe-inline'; "
        "img-src 'self' data: blob:; "
        "font-src 'self' data:; "
        "connect-src 'self'; "
        "base-uri 'self'; form-action 'self'"
    )
    return resp

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=int(os.environ.get("PORT", 8085)))
