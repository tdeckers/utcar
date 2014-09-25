# UTC Alarm Receiver

utcar acts as a central station to a ATS alarm system.

My use case:
* leverage the alarms systems events (e.g. motion, window sensors) as input to home automation

`utcar` listens by default on port number 12300, can be set on command line (`-port`)

The alarm must be configured to send OH+XSIA messages. In that case it will send two types of messages: heartbeats and (X)SIA messages.
The heartbeats look like this:

`SR0001L0001    001465XX    [ID5B9490D8]`

At this point, heartbeats are logged but ignored further.

The (X)SIA messages are more interesting, the look like this:

`01010053"SIA-DCS"0007R0073L0011[#001365|NUA021*'detector hall'NM]7C9677F21948CC12|#001365`

This is a message to indicate activation (UA) of a motion sensor in zone 21 (detector in hall).  If no `-taddr` is provided, we only log this message.

if a `-taddr` parameter is provided, then `utcar` will POST a message to an HTTP endpoint. Right now, this is customized for Openhab - might need to generalize this later.

URL for the POST: `https://<taddr>/rest/items/al_{item}/state`, where item is the zone received from the alarm.

For the example above, an HTTP POST with body ON is sent to:

`https://<taddr>/rest/items/al_021/state`

If an (X)SIA message of UR is received, an OFF message is sent. Support for more messages types might be added later.
