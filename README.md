# Service-in-go-
the project is an API rest that show the server's information of a given domain.the application also show the domains hat recently has 
searched in the application 

# Server information 
the information that the endpint collect is:
- servers:contain an servers array asociates to the domain, each object of the array contains:
    - address: Ia IP or the host of the server.
  - ssl_grade: the sslgrade calificated to SSLabs
  - country: the country of the IP, this information was obtain usging WHOIS API 
  - owner: the organization owner of the Ip, this information was obtain usging WHOIS API too.
- servers_changed:The go aplication calculates if the servers has changed since one hour, only if the domain was saved in the application before.
this parameter will return true if the server has changed since one hour of before 
- ssl_grade: this is the lower ssl-grade of all the servers.
- previous_ssl_grade: this is the lower ssl-grade that the servers had since a hour before.
- logo: this show the domain logo if this exists, this logo is obtain from the API besticon-demo
- title: this show the title of the page obtain of the <head> of the domain.
- is_down: this return true if the server is down.

# Run the application.
this is a local aplication so its necesary let running the "main.go" element,  this will run on localhost and the port 3000,
and the appication an be used fom the index element.

# Technologies used.
the was created using Golang to create the endpoint and consume diferents APIS to collect the servers information, VueJs to show 
the elements inside the page and postgres to stock the servers information.

