echo {> config.json
echo "api": "%API_ENDPOINT%",>> config.json
echo "apps_domain": "%APPS_DOMAIN%",>> config.json
echo "admin_user": "%ADMIN_USER%",>> config.json
echo "admin_password": "%ADMIN_PASSWORD%",>> config.json
echo "cf_user": "%CF_USER%",>> config.json
echo "cf_user_password": "%CF_USER_PASSWORD%",>> config.json
echo "cf_org": "%CF_ORG%",>> config.json
echo "cf_space": "%CF_SPACE%",>> config.json
echo "skip_ssl_validation": true,>> config.json
echo "persistent_app_host": "persistent-app-win64",>> config.json
echo "default_timeout": 120,>> config.json
echo "cf_push_timeout": 210,>> config.json
echo "long_curl_timeout": 210,>> config.json
echo "broker_start_timeout": 330>> config.json
echo }>> config.json