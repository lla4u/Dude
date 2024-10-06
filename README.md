# Dude - Dynon user data explorer

> Each SkyView HDX display can act as a Primary Flight Display (PFD) with Synthetic Vision, an
Engine Monitoring System (EMS), and a Moving Map in a variety of customizable screen
layouts. Data is sourced from various connected modules and devices.
>
> SkyView HDX displays record and store flight information in several datalogs which can be
exported for analysis by the owner, and a high-resolution datalog which can be used by Dynon
for troubleshooting. 
>
> This tool provide an easy and efficient way to enlighten the Dynon datalogs for:  
>  - Providing long term flight history,
>  - Providing capability to display flight map and related parameters for improving pilote usage and safety,
>  - Providing accurate informations like the average speed to considere preparing flight, the average landing speed,
>  - ...

## What are the dependencies
> This tool use:
> - Docker for containers and network management
> - Influx database for timed long term storage
> - Gafana for data presentation
> - Utility program (Mac, Windows, Linux) for onboarding Dynon datalogs into database.
>
> All of the selected are open source and free for personal usage.

## What is provided
> A local web interface to query flight data shuch as:
> ![Screenshot of web interface.](https://github.com/lla4u/Dude/blob/main/Screenshots/Screenshot_web_Interface.png)
>
> ![Screenshot of web interface 2.](https://github.com/lla4u/Dude/blob/main/Screenshots/Screenshot_web_interface_2.png)

# Installation procedure


## 1 - Docker install (if not yet done)
Depending your operating system download and install docker software.
https://docs.docker.com/get-docker/  

Video : https://www.youtube.com/watch?v=mS26N5cLBe8&ab_channel=CodersArcade  


## 2 - Building Dude stack
```
1. Where the stack will be sitting:
   Open terminal, cmd and create your home install directory 
     cd /home/lla 
   then move into: 
     cd /home/lla

2. Clone de github repo:
   git clone https://github.com/lla4u/Dude.git
   or
   Download and unzip zip archive downloaded from github.

3. change directory to Dude-Influx-Grafana
   cd Dude

4. Build the Docker stack using terminal: 
   docker-compose --env-file config.env up --build -d 
   or (for recent docker version) 
   docker compose --env-file config.env up --build -d 

   after a while (mostly depending your network bandwith) 3 containers will be created and available.

5. Check:
   execute docker ps from the terminal

   Having 2 conainers running you are good to go further ...
```
> [!IMPORTANT]
> Persistant data (InfluxDB & Grafana) will be stored into a local directory (Docker). Removing will result a loss of data.

## 2 - Starting / Stoping Dude stack
> Starting or stoping the Dude stack can be acheived using:
> 1. Docker Dashboard
>    - Starting
> ![Screenshot Docker dashboard starting.](https://github.com/lla4u/Dude/blob/main/Screenshots/Screenshot_docker_dashboard_start.png)
>
>    - Stoping
> ![Screenshot Docker dashboard stoping.](https://github.com/lla4u/Dude/blob/main/Screenshots/Screenshot_docker_dashboard_stop.png)
>
> 2. Command line
>    - Open terminal and move to the Dude-Influx-Grafana directory
>      - cd /home/lla/dude/Dude
>
>    - Starting
>      - docker compose up
>
>    - Stoping
>      - docker compose down

# Adding HDX datalogs to the solution
> Adding datalogs in the tool is a two steps process:
> - First, Collect datalog from the HDX
> - Second, Import the datalogs into Influxdb using Dude utility program.

## Collecting datalog from the HDX
> Collecting datalog from the HDX is quite trivial and require usb key plugged into the Dynon:
> ( I use the same usb key that for the plates and map updates ...)
> 1. Fire up your HDX
> 
> 2. Press button 7 & 8 Simultaneously for few seconds to startup the dynon setup screen
> 
> 3. Navigate to SYSTEM SOFTWARE -> EXPORT USER DATA LOGS 
> 
> 4. (Otional) Define label
> 
> 5. Export pressing button 8
> 
> Video: https://www.youtube.com/watch?v=fS6H_8gNd90&ab_channel=RobertHamilton

> [!IMPORTANT]
> Dynon datalog storage is limited and rewrited over time. So collect datalogs around every 8 flight hours or accept to loose information.

## Importing collected datalog into dude stack
> The docker stack is provided with an additional utility program that is available in the Go/build directory.  
> There is one version per common operating system such as Darwing, Windows, Linux.
> You can copy and rename the file in existing PATH or add the directory in PATH to get it available from anywhere.
>
> From teminal or cmd or powershell (windows) execute:  
>  Dude_darwin_amd64 import --help (for Intel Macbook)  
>  Dude_windows_amd64 import --help (for Intel windows)  
>  ...  
>  This Will give you the expected parameters if any.  

> [!TIP]
>  It is recomended to create an hiden file (.Dude.yaml) in your $HOME directory to bypass the parameters.  

>  The following file must have the **datalog parameter adjusted to your datalogs directory path!**.
> 
```
❯ cat $HOME/.Dude.yaml
---
# Specify the influx database url
url: http://localhost:8086

# Specify the influx database token
token: my-super-secret-auth-token

# Specify the directory path of your datalogs
datalog: /Users/lla/Documents/Laurent/Aviation/P300 Dude
```

> ![Screenshot of Dude help.](https://github.com/lla4u/Dude/blob/main/Screenshots/Screenshot_dude_cli_help.png)
>
> Process is as follow:
> 1. Copy the usb key csv file(s) (USER_DATA_LOG.csv) into the Datalogs directory.
>
> 2. Execute (Intel Mac):
>  ./dudeImport_darwin_amd64 import
> 
> ![Screenshot of Dude.](https://github.com/lla4u/Dude/blob/main/Screenshots/Screenshot_dude_cli.png)
> Screenshot note:  
> - Onboarding datalogs required 1.63 seconds  
> - Submited datalogs file was having 164897 csv rows  
> - Import saved 28860 rows into influx database.  
>
> Using datalog parameter:
> ![Screenshot of dudeImport.](https://github.com/lla4u/Dude/blob/main/imageScreenshots/Screenshot_dude_cli_datalog.png)

> [!NOTE]
> Dude Import utility tool only import new datalogs! Imported.txt file is created into the datalog directory after import and hold the already imported datalog files.

# Roadmap
- [ ] Use Grafana variable to help finding the flights saved into the InfluxDB.
- [ ] Create document for helping users to use the Grafana UI and look at the Dynon datalogs.

Have safe flights.  

Laurent
