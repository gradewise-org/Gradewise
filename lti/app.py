import os, json
from flask import Flask, request, jsonify, redirect, session
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
app.config["CACHE_TYPE"] = "SimpleCache"
app_cache = Cache(app)

launch_store = FlaskCacheDataStorage(app_cache)

# -- ENDPOINTS -- #
@app.route("/lti/jwks.json")
def jwks():
    # Serve public keys so Canvas can validate any JWTs (e.g., for Deep Linking responses).
    return jsonify(TOOL_JWKS)



@app.route("/lti/oidc-login", methods=["GET", "POST"])
def oidc_login():
    req = FlaskRequest()
    try:
        # Log what Canvas actually sent (works for GET or POST)
        print("OIDC iss:", req.get_param("iss"), " client_id:", req.get_param("client_id"))

        login = FlaskOIDCLogin(
            req,
            tool_conf,
            launch_data_storage=launch_store
        ).enable_check_cookies()

        target_link_uri = req.get_param("target_link_uri")
        if not target_link_uri:
            # Canvas must send this; bail early with detail if missing
            return jsonify({"error": "missing target_link_uri"}), 400

        # Always redirect (both on POST and on the GET preflight)
        return login.redirect(target_link_uri)

    except Exception as e:
        # Show the *real* reason for the 400 you’re seeing
        print("OIDC ERROR:", type(e).__name__, str(e))
        return jsonify({"status": "error", "error": type(e).__name__, "message": str(e)}), 400



@app.route("/lti/launch", methods=["POST"])
def launch():
    req = FlaskRequest()
    try:
        ml = FlaskMessageLaunch(req, tool_conf, launch_data_storage=launch_store).validate()
    except Exception as e:
        app.logger.exception("LTI launch validate failed")
        return jsonify({"status": "error", "message": str(e)}), 400

    launch_data = ml.get_launch_data()
    session["lti_sub"] = launch_data.get("sub")
    session["lti_context"] = launch_data.get("https://purl.imsglobal.org/spec/lti/claim/context")
    resp = redirect("/app/")
    resp.set_cookie(
    "launched", "1",
    path="/",
    secure=True,
    samesite="None",
    max_age=3600
    )
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
