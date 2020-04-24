# UTC Alarm Receiver

_utcar_ acts as a central station to a ATS alarm system and has been tested using an ATS2000IP system.  It might work with other variants as well.  The communication between the alarm system and utcar is using OH+XSIA protocol.  I'm not an expert, but I assume that any system that uses this protocol should be compatible (naive? me?).

My use case:
* leverage the alarms systems events (e.g. motion, window sensors) as input to home automation

## Configure your alarm system

Your alarm is typically configured to send messages to a central station in case of alarm, fire, etc...  The link is checked for health by sending periodic heartbeat messages.  In my case, these heartbeat messages are sent every minute.

For _utcar_ to work, we'll configure it in the alarm system as a new central station.  Since we don't want to interrupt messages going to the first central station (the one calling the police in case of an alarm), we need to go through a few steps to make this work.  Let's assume we want to get notified when a motion sensor is triggered.

1. Create a new area that we'll use with _utcar_.  This allows us to keep automation for _utcar_ separate from the actual alarms going to the formal central station.  Let's pick area 4.
  ![alt text](https://github.com/tdeckers/utcar/raw/master/img/area.png)
2. Configure a filter that follows the motion detector
  ![alt text](https://github.com/tdeckers/utcar/raw/master/img/filter.png)
3. Configure an output to follow the filter
  ![alt text](https://github.com/tdeckers/utcar/raw/master/img/output.png)
4. Configure a new zone, and configure the output as its _virtual zone_.  Configure it in area 4.
  ![alt text](https://github.com/tdeckers/utcar/raw/master/img/zone.png)
5. Configure a new central station (IP based), configure it with the host IP address of the machine where you'll be running _utcar_.  Also pick a port number, by default _utcar_ runs on port 12300.  Associate this central station with area 4.
  ![alt text](https://github.com/tdeckers/utcar/raw/master/img/central_station.png)


## Running utcar

	Usage of ./utcar:
	  -port=12300: Listen port number (default: 12300)
	  -taddr="": Target addr (host:port)
	  -tpwd="": Target password
	  -tuser="": Target username

Example:

	./utcar -port=10000 -taddr=localhost:8443 -tuser=yourname -tpwd=yourpass
	2014/09/25 06:40:32 Listing on port 10000...
	2014/09/25 06:40:32 Pushing to localhost:8443

_utcar_ listens by default on port number 12300, can be set on command line (`-port`)

The alarm must be configured to send OH+XSIA messages. In that case it will send two types of messages: heartbeats and (X)SIA messages.
The heartbeats look like this:

	SR0001L0001    001465XX    [ID5B9490D8]

At this point, heartbeats are logged but ignored further.

The (X)SIA messages are more interesting, the look like this:

	01010053"SIA-DCS"0007R0073L0011[#001365|NUA021*'detector hall'NM]7C9677F21948CC12|#001365

This is a message to indicate activation (UA) of a motion sensor in zone 21 (detector in hall).  If no `-taddr` is provided, we only log this message.

if a `-taddr` parameter is provided, then _utcar_ will POST a message to an HTTP endpoint. Right now, this is customized for Openhab - might need to generalize this later.

URL for the POST: `https://<taddr>/rest/items/al_{item}/state`, where item is the zone received from the alarm.

For the example above, an HTTP POST with body ON is sent to:

	https://<taddr>/rest/items/al_021/state

You can provide `-tuser` and `-tpwd` to provide basic authentication.

If an (X)SIA message of UR is received, an OFF message is sent. Support for more messages types might be added later.

## Running in a container

You can also run utcar in a container:

	docker run -p 12300:12300 tdeckers/utcar -port=10000 -taddr=localhost:8443

# Building

To cross-compile for different platforms, set up your cross compilers for the platforms you need.

	cd $GOROOT/src
	# Windows 64-bit
	GOOS=windows GOARCH=amd64 ./make.bash --no-clean
	# Raspberry Pi
	GOOS=linux GOARCH=arm ./make.bash --no-clean

Then to compile your app:

	cd $YOUR_APP
	GOOS=windows GOARCH=amd64 go build
	
Note: 3 versions are provided: for Windows 64, arm32 (for raspberry pi) and for linux64. See releases.

## Building a container
(thanks: https://rollout.io/blog/building-minimal-docker-containers-for-go-applications/)

First build a statically linked executable:

	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -installsuffix cgo
	cp /etc/ssl/certs/ca-certificates.crt .
	docker build -t utcar .

Then run the container:

	docker run -p 12300:12300 utcar /utcar -port=10000 -taddr=localhost:8443

Credits: Thanks to Dirk @ OP for his help on the ATS configuration.
