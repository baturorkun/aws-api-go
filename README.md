# AWS-API

**It's an example of a API using Aws-Sdk-Go (GoLang) with Gorilla Toolkit**

Aws-sdk-go is the official AWS SDK for the Go programming language.

##Used libraries from GitHub

- Gorilla Web Toolkit
- Aws-SDK

## Requirements
- go 1.13
- export GO111MODULE=on

`code()`
## Usages

#### Running on Local Machine

`$ go run mail.go`

#### Running on Remote Machine

`$ sh deploy.sh batur baturorkun.com`

#### Running on Docker

```bash
$ docker build -f docker/Dockertfile --build-arg PROJECT="aws-api" --build-arg USER_ID=`id -u` -t aws-api .
$ docker run --env-file=.env -v $PWD:/builder/src/aws-api --rm aws-api
```

#### Running on Docker-Compose

`$ docker-compose up`

### API Methods

- #### /instance-search
    - GET Params: 
     
        > tag_name : String , Ex: Name
                         
        > tag_value : String
        
- #### /instance-start 
    - GET Params: 
        
         > instance-id : String , Ex: i-30ffAddYT65Ahj890
        
- #### /instance-stop
    - GET Params:
    
        > instance-id : String , Ex: i-30ffAddYT65Ahj890

- #### /instance-terminate
    - GET Params:
    
        > instance-id : String , Ex: i-30ffAddYT65Ahj890
         
- #### /instance-setting-getdisabletermination
  - GET Params:
    
        > instance-id : String , Ex: i-30ffAddYT65Ahj890

- /instance-setting-setdisabletermination       
    - GET Params:
    
        > instance-id : String , Ex: i-30ffAddYT65Ahj890           
                           
- #### /snapshot-create
    - GET Params:
    
        > instance-id : String , Ex: i-30ffAddYT65Ahj890

        > tag-name : String , Ex: Name
        
        > tag-value : String , Ex: Batur                                                                              

        > state-name: String , Ex: running or stopped
                                                                                                                                                                                                                                                                                                 
        > public-ip : String , Ex: 192.168.1.100  

- #### /elasticip-allocate
    - GET Params:
    
        > number : Integer , Ex: 1,2,3,... (How many IPs)

- #### /elasticip-search
    - GET Params:
    
        > association-id : String ; Optional , Ex: eipalloc-0b768f070efba5132    

- #### /release-elasticip
    - Not Ready
    
- #### /lbtargetgroup-search
    - GET Params:
    
        > names : String ; Coma seperated string , Ex: test1,test2,test3

- #### /tag-delete
    - GET Params:
    
        > instance-id : String , Ex: i-30ffAddYT65Ahj890

        > tag : String , Ex: Name
             
        > value : String , Ex: Batur     

- #### /tag-create
    - GET Params:
    
        > instance-id : String , Ex: i-30ffAddYT65Ahj890

        > tag : String , Ex: Name
             
        > value : String , Ex: Batur   
                                                                                                                                                     
- #### /billing-daily
    - No Parameters
                                                                                                                                                         
- #### /billing-monthly
    - No Parameters   
                                                                                                                                     
- #### /remote-copy-sshkey
    - GET Params:
    
        > ssh-key : String ; Filename , Ex: batur.pem

        > public-ip : String , Ex: 192.168.1.100
                                                                                                                                       
- #### /remote-get-messsages-log
     - GET Params:
     
        > lines : Integer ; Lines Number , Ex: 100

        > public-ip : String , Ex: 192.168.1.100
                                                    
- #### /sendemail
     - GET Params:
     
        > file : String ; Filename; Optinal , Ex: file.dat

        > recipients : String ; Coma seperated E-mail addresses , Ex: batur@domain.com,orkun@domain.com

        > subject : String , Email subject
 
        > body : String , Email content
          
                                                                                                                                                                                                                                                
#### License

Copyright © 2020 Batur Orkun
Distributed under the Eclipse Public License either version 1.0 or (at your option) any later version.
