import os, json
from flask import Flask, request, jsonify
from dotenv import load_dotenv
from pylti1p3.contrib.flask import (
    FlaskOIDCLogin, FlaskMessageLaunch, FlaskRequest, FlaskCacheDataStorage
)
from pylti1p3.tool_config import ToolConfDict

load_dotenv("config/.env")

# ---- Tool Configuration ----- #
ISSUER = "https://canvas.instructure.com"

tool_conf = ToolConfDict({
    ISSUER: [{
        "default": True,
        "client_id": os.environ["CANVAS_CLIENT_ID"],
        "auth_login_url": os.environ["CANVAS_OIDC_AUTH_URL"],
        "auth_token_url": os.environ["CANVAS_TOKEN_URL"],
        "auth_audience": os.environ.get("CANVAS_TOKEN_AUDIENCE"),  # often same as token URL; ok if missing
        "key_set_url": os.environ["CANVAS_JWKS_URL"],
        # For deep linking / tool-signed JWTs, set one of the following:
        # "private_key_file": "config/private.key",
        # "public_key_file": "config/public.key",
        # or inline:
        # "private_key": os.environ.get("TOOL_PRIVATE_KEY_PEM"),
        # "public_key": os.environ.get("TOOL_PUBLIC_KEY_PEM"),
        "deployment_ids": [os.environ.get("CANVAS_DEPLOYMENT_ID", "")]  # optional but recommended
    }]
})


# ---- Flask App ---- #
app = Flask(__name__)
app.secret_key = os.environ.get("FLASK_SECRET", "change-me")
TOOL_JWKS = {"keys": []}  # replace with  actual tool JWKS later when we start signing responses


@app.route("/jwks.json")
def jwks():
    # Serve public keys so Canvas can validate any JWTs (e.g., for Deep Linking responses).
    return jsonify(TOOL_JWKS)

@app.route("/oidc-login", methods=["POST"])
def oidc_login():
    # Canvas posts iss, login_hint, target_link_uri (and maybe lti_message_hint)
    login = FlaskOIDCLogin(FlaskRequest())
    iss = request.form["iss"]
    login_hint = request.form.get("login_hint", "")
    target_link_uri = request.form["target_link_uri"]
    lti_message_hint = request.form.get("lti_message_hint")
    return login.redirect(
        tool_conf, iss, login_hint,
        target_link_uri=target_link_uri,
        lti_message_hint=lti_message_hint
    )

@app.route("/launch", methods=["POST"])
def launch():
    # Validate id_token (signature, iss/aud/exp, state/nonce)
    ml = FlaskMessageLaunch(FlaskRequest(), tool_conf).validate()
    launch_data = ml.get_launch_data()
    # TODO: persist ml/launch_data for later AGS calls: ml.get_ags(), ml.get_nrps()
    # For now, just confirm it worked.
    return jsonify({
        "status": "ok",
        "platform": launch_data.get("iss"),
        "user_sub": launch_data.get("sub"),
        "roles": launch_data.get("https://purl.imsglobal.org/spec/lti/claim/roles", []),
        "context": launch_data.get("https://purl.imsglobal.org/spec/lti/claim/context", {})
    })

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=int(os.environ.get("PORT", 8080)))
