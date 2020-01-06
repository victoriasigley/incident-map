# incident-map
Enrich 911 emergency incident data to provide better analytics for a fire department.

Steps to install and run application:
- clone repository
- Prepare sql:
  - make sure pg_hba.conf connections are set to trust (can find path by entering postgres command line prompt, with command 'show hba_file;')
  - cd sql
  - psql -U postgres -f 0001_init_incident.up.sql
- Prepare server:
  - cd incident-map
  - source dev/incident-map.env
  - go run server/server.go &
- cd incident-map && npm start
- access in the browser at <your-ip>:3000 OR view live hosted on my AWS instance http://3.226.91.12:3000/

Future improvements?
- Use martini (or other router) to handle server side routing
- Upgrade postgres on my machine to enable using ON CONFLICT commands (see commented out code in 0001_init_incident.up.sql)
- Implement migrate in order to enable server version control and handling for future sql files
- Imporved map markers (would need to be custom in order to use different colors/icons). Would implement this such that each apparatus had its own separate color, possibly a label with either the name of the apparatus or an icon. Would be nice to have the incident marker popup open on start-up, but was not straight-forward to implement (requires custom components).
-Ability to change files for upload (without having to cancel out and click upload again).

Time spent on project: ~10 hours
