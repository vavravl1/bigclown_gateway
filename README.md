# bigclown_gateway
Go implementation of Bigclown CoreModule->(raspberry/omnia) gateway

## Required tools
* make
* go (go version go1.10.2)
* dep (v0.4.1)

## How to clone
* git clone to your $GOPATH/src/github.com/vavravl1/bigclown_gateway
* cd $GOPATH/src/github.com/vavravl1/bigclown_gateway

## How to build & deploy
### Before first deployment
* ssh to your node (raspberry/turris)
* create folder for the gateway, e.g. ```mkdir bigclown_gateway```
* on your localhost
  * ```scp $GOPATH/src/github.com/vavravl1/bigclown_gateway/bigclown_gateway.sh <your-node>:/home/<user-name>/bigclown_gateway```
* On the node 
  * Update variables in ``/home/<user-name>/bigclown_gateway/bigclown_gateway.sh``
  * Upate crontab to run the gateway after reboot 
    * ```crontab -e ```
    * Add on the end of the file ```@reboot /home/<user-name>/bigclown_gateway/bigclown_gateway.sh```
  * Reboot, e.g. ```sudo reboot```
  * Follow the steps below to deploy the binary

### Regular deployments
* ```cd $GOPATH/src/github.com/vavravl1/bigclown_gateway```
* Update the deploy part of the makefile to address your raspberry settings
* ```make```
