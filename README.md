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
>

## What are the dependencies
> This tool use:
> - Docker for containers and network management
> - Influxdb  database for timed long term storage
> - Gafana for data presentation
>
> All of the selected are open source and free for personal usage.
>

## What is provided
> - A local web interface to query flight data shuch as:
> ![Screenshot of web interface.](https://github.com/lla4u/Dude/blob/main/Screenshots/Screenshot_web_Interface.png)
>
> - A dude utility program (MacOs, Windows, Linux) for onboarding and visualyzing Dynon datalogs information.
>

## What is required
> Internet connectivity
>

# Installation procedure

## 1 - Docker install (if not yet done)
> Depending your operating system download and install docker software.
> https://docs.docker.com/get-docker/  
> 
> Video : https://www.youtube.com/watch?v=mS26N5cLBe8&ab_channel=CodersArcade  
> 

## 2 - Building Dude stack

### 2.1 Where the stack will be sitting:
>    Open terminal, cmd and move into your HOME directory 
>
> ```shell
> cd $HOME 
> ```
>

### 2.2. Clone or download and unzip the github repo into HOME directory:
>
> ```shell
> git clone https://github.com/lla4u/Dude.git
> ```
>
> or
>
> Download and unzip downloaded compress archive into $HOME
>
> ```shell
> https://github.com/lla4u/Dude/archive/refs/heads/main.zip
> ```
>

### 2.3. Change current directory to Dude
> 
> ```shell
> cd Dude
> ```
>

### 2.4. Build the Docker stack using terminal: 
>
> ```shell
> docker compose --env-file config.env up --build -d 
> ```
>
> or (for older docker version) 
>
> ```shell
> docker-compose --env-file config.env up --build -d 
> ```

> [!NOTE]
> Wait until the 2 containers downloaded and running. Mostly depending your network bandwith. 

### 2.5. Check using terminal:
> 
> ```shell
> docker ps 
> ```
> Having 2 containers running you are good to go further ...
>

> [!IMPORTANT]
> Dude Persistant data (InfluxDB & Grafana) will be stored into a local directory (Docker). Removing this directory will result a loss of data.

### 2.6. Setup the Datalogs directory
> Dude cli tools require to know where to find your collected HDX datalogs.
> The file that hold this configuration must be sitting into your HOME directory with .Dude (hiden file) name and yaml as extension.
> You must adapt the provided for making:
> 1. the DatalogPath value to link with your datalogs store directory.
> 2. the Location value relevant to your time zone. (Mandatory to align stats to your Location/time zone)
> 
> ```shell
> echo '---
> # Specify the directory path of your datalogs
> # Must be adapted to your own preference and located in your HOME directory !!!
> # Such as: .Dude.yaml
> DatalogPath: "/Users/lla/Documents/Laurent/Aviation/P300 Dude"
> Location: "Europe/Paris"' > $HOME/.Dude.yaml
>```
> 

# Starting / Stoping Dude stack
> Starting or stoping the Dude stack can be acheived using:
> 

## 1. Docker Dashboard
>    - Starting
> ![Screenshot Docker dashboard starting.](https://github.com/lla4u/Dude/blob/main/Screenshots/Screenshot_docker_dashboard_start.png)
>
>    - Stoping
> ![Screenshot Docker dashboard stoping.](https://github.com/lla4u/Dude/blob/main/Screenshots/Screenshot_docker_dashboard_stop.png)
>

## 2. Command line
> Open terminal and move to the Dude directory
> ```shell
> cd $HOME/Dude
> ```
>
> - Starting
> ```shell
> docker compose up
> ```
> - Stopping
> ```shell
> docker compose down
>

# Adding HDX datalogs to the solution
> Adding datalogs in the tool is a two steps process:
> - First, Collect datalog from the HDX
> - Second, Import the datalogs into Influxdb using Dude utility program.
> 

## Collecting datalog from the HDX
> Collecting datalog from the HDX is quite trivial and require usb key plugged into the Dynon:
> ( I use the same usb key that for the plates and map updates ...)
>
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
> The Dude stack is provided with an additional utility program available in the Go/build directory.  
> There is one version per common operating system such as Darwing, Windows, Linux.
> You can copy and rename the file in existing PATH or add the build directory in PATH to get it available from anywhere.
>
> Using terminal or cmd / powershell (windows) execute:  
> 
> Dude_darwin_amd64 --help (for Intel Macbook)
>   
> Dude_windows_amd64 --verbose (for Intel windows)  
>
>  ...
>  
>  flags --help or --verbose will give you the expected parameters if any.  

> ![Screenshot of Dude help.](https://github.com/lla4u/Dude/blob/main/Screenshots/Screenshot_dude_cli_help.png)

> Import process in Dude stack is as follow:
> 1. Copy the usb key csv file(s) (USER_DATA_LOG.csv) into your Datalogs directory.
>
> 2. Execute (Intel Mac):
>  ./dudeImport_darwin_amd64 import
> 
> ![Screenshot of Dude.](https://github.com/lla4u/Dude/blob/main/Screenshots/Screenshot_dude_cli.png)
> Screenshot note:  
> - Onboarding datalogs required 1.42 seconds.
> - Submited datalogs file was having 164.897 csv rows.
> - Import saved 9805 rows into influx database.
> - Import skiped 76042 rows as flight(s) was already known.  
>

> [!NOTE]
> Dude import utility tool only import new datalogs as an Imported.yml file is created into the datalog directory. This important file hold all the known Datalogs and relevant Flights information. 

## Viewing collected datalog and flights

> Dude stats utility expose a sumary and flights details.
> 
>```shell
> ./Dude_darwin_amd64 stats
> ```
> 
> ![Screenshot of Dude stats.](https://github.com/lla4u/Dude/blob/main/Screenshots/Screenshot_dude_cli_stats.png)
> Screenshot note: 
> - Sumary of Datalogs, relevant Flights and flights duration are presented
> - Each flight and relevant start / End / Duration is listed as well as a link for you to copy paste in your favorite internet browser.
> (you might have to sign in grafana (admin/admin) the first time)
>

# Roadmap
- [ ] Create document for helping to use Grafana UI and tune up presentation at your taste.

Have safe flights.  

Laurent
