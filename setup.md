```sh
# Create admin user
# Run on mattermost container
/mattermost/bin/mmctl --local user create --email admin@example.org --username admin --password password --system-admin

# Deploy
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_USERNAME=admin
export MM_ADMIN_PASSWORD=password

make deploy
 ```