# incident-map
Enrich 911 emergency incident data to provide better analytics for a fire department.

Steps to install and run application:
- clone repository
- cd incident-map
- source dev/incident-map.env
- go run server/server.go &
- cd incident-map && npm start
- access in the browser at <your-ip>:3000 OR view live hosted on my AWS instance http://3.226.91.12:3000/

What improvements would I make with double the time?
- I didn't have time to finish my server side implementation. I need to handle posts to my endpoint and store in a database. With double the time, I believe I would've had time to complete this. I am also more familiar with using martini to handle my endpoints and may have used that, but determined I didn't have enough time and was relying on the go included libraries.
- I didn't have time to create my postgreSQL database. On a post to my server's incident endpoint, I would've liked to store all the data fields in the database (or at least the ones I cared about- my server struct I did all to start out with because I hadn't determined 100% which of the data I would be most interested in using, I figure the more info the better). I have postgres set up on my machine, but typing out all the field names takes a long time. I could've used the stringified and encoded/decoded but prefer not doing this. Given twice the time, I believe I would have had time to complete this as well.
- The data output I was intending on using for the front end wasn't very sophisticated, just taking all the lat/lon points and displaying them on the map. With more time, I would've liked to search through the unit_status and either just displayed the most recent update, or display them in a way where older statuses were possibly more transparent to indicate they were an earlier data point. Below the map I also would have liked to include more information from the description (comments, type, subtype, first_unit_arrived, first_unit_enroute, first_unit_dispatched for starters), address summary, etc. I also would have liked to add a map marker with the incident latitude/longitude with a distinct marker to indicate that it's the origin of the incident. I believe a lot of these sophistications (or maybe all) could have been completed given twice the time.

Time spent on project: 4 hours
