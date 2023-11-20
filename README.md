# dfxp
Dflux Performance 

## ZMQ overview 
messages list:
    1.  start request/response 
    2.  stop request/response
    3.  get info request/response
    4.  tunnels add 
    5.  tunnels delete
    6.  get info     
    7.  publish metrics
    8. error response 
    9. msg error response 
    
## ZMQ messages format
### 1 Start message
    1. Request 
        |  byte4 |  byte3 |  byte2 |  byte1 | 
        | ------ | ------ | ------ | ------ |
        | command         | length          | 
        | ------ | ------ | ------ | ------ | 
        |             flow id               |
        | ------ | ------ | ------ | ------ |
        |     metric interval (sec)         |
        | ------ | ------ | ------ | ------ |                  

    2. Response
        |  byte4 |  byte3 |  byte2 |  byte1 | 
        | ------ | ------ | ------ | ------ |
        | command         | length          |
        | ------ | ------ | ------ | ------ |
        |             flow id               |
        | ------ | ------ | ------ | ------ |
        | char4  |  char3 |  char2 |  char1 | 
        | ------ | ------ | ------ | ------ |
        | charN  |      .....               |  
        | ------ | ------ | ------ | ------ |
        
        where N - max publisher length

### 2 Stop message
    1. Request 
        |  byte4 |  byte3 |  byte2 |  byte1 | 
        | ------ | ------ | ------ | ------ |
        | command         | length          | 
        | ------ | ------ | ------ | ------ | 
        |             flow id               |
        | ------ | ------ | ------ | ------ |

    2. Response
        |  byte4 |  byte3 |  byte2 |  byte1 | 
        | ------ | ------ | ------ | ------ |
        | command         | length          |
        | ------ | ------ | ------ | ------ |
        |             flow id               |
        | ------ | ------ | ------ | ------ |
       

### 3 Info message
    1. Request 
        |  byte4 |  byte3 |  byte2 |  byte1 | 
        | ------ | ------ | ------ | ------ |
        | command         | length          | 
        | ------ | ------ | ------ | ------ | 
        |             flow id               |
        | ------ | ------ | ------ | ------ |
     
    2. Response
        |  byte4 |  byte3 |  byte2 |  byte1 | 
        | ------ | ------ | ------ | ------ |
        | command         | length          |
        | ------ | ------ | ------ | ------ |
        |             flow id               |
        | ------ | ------ | ------ | ------ |\
        | char4  |  char3 |  char2 |  char1 | |
        | ------ | ------ | ------ | ------ | |  version[DFXP_ZMQ_MAX_VER_LEN = 16]
        | charN  |      .....               | | 
        | ------ | ------ | ------ | ------ |/
        
        where N - max info response length
        info context:
            1. version (max length) 

### 4 Add tunnels message
    1. Request 
        |  byte4 |  byte3 |  byte2 |  byte1 | 
        | ------ | ------ | ------ | ------ |
        | command         | length          | 
        | ------ | ------ | ------ | ------ | 
        |             flow id               |
        | ------ | ------ | ------ | ------ |
        |             tunnels number        |
        | ------ | ------ | ------ | ------ |
        |             tunnel1               |
        | ------ | ------ | ------ | ------ |
        |           ...........             |
        | ------ | ------ | ------ | ------ |
        |             tunnelN               |
        | ------ | ------ | ------ | ------ |


        - Tunnel 
            |  byte4 |  byte3 |  byte2 |  byte1 | 
            | ------ | ------ | ------ | ------ |
            |                 teid in           | 
            | ------ | ------ | ------ | ------ | 
            |                 teid out          | 
            | ------ | ------ | ------ | ------ |
            |                 ue_ipv4           | 
            | ------ | ------ | ------ | ------ |
            |                 upf_ipv4          | 
            | ------ | ------ | ------ | ------ |


    2. Response
        |  byte4 |  byte3 |  byte2 |  byte1 | 
        | ------ | ------ | ------ | ------ |
        | command         | length          |
        | ------ | ------ | ------ | ------ |
        |             flow id               |
        | ------ | ------ | ------ | ------ |
        |             tunnels number        |
        | ------ | ------ | ------ | ------ |

### 5. Remove tunnels message

    1. Request 
        |  byte4 |  byte3 |  byte2 |  byte1 | 
        | ------ | ------ | ------ | ------ |
        | command         | length          | 
        | ------ | ------ | ------ | ------ | 
        |             flow id               |
        | ------ | ------ | ------ | ------ |
        |             teid number           |
        | ------ | ------ | ------ | ------ |
        |             teid1                 |
        | ------ | ------ | ------ | ------ |
        |           ...........             |
        | ------ | ------ | ------ | ------ |
        |             teidN                 |
        | ------ | ------ | ------ | ------ |



    2. Response
        |  byte4 |  byte3 |  byte2 |  byte1 | 
        | ------ | ------ | ------ | ------ |
        | command         | length          |
        | ------ | ------ | ------ | ------ |
        |             flow id               |
        | ------ | ------ | ------ | ------ |
        |             tunnels number        |
        | ------ | ------ | ------ | ------ |


### 6. Metrics message

 1. Metric publish message

        |  byte4 |  byte3 |  byte2 |  byte1 | 
        | ------ | ------ | ------ | ------ |
        | command         | length          | 
        | ------ | ------ | ------ | ------ | 
        |             flow id               |
        | ------ | ------ | ------ | ------ |
        |             metrics number        |
        | ------ | ------ | ------ | ------ |
        |             metric1               |
        | ------ | ------ | ------ | ------ |
        |           ...........             |
        | ------ | ------ | ------ | ------ |
        |             metricN               |
        | ------ | ------ | ------ | ------ |

        - Metrics:

            |  byte4 |  byte3 |  byte2 |  byte1 | 
            | ------ | ------ | ------ | ------ |
            | command         | length          | 
            | ------ | ------ | ------ | ------ | 
            |             flow id               |
            | ------ | ------ | ------ | ------ | 
            |                          |protocol|         
            | ------ | ------ | ------ | ------ | 
            |     UDP metrics  or               |
            |     TCP metrics  or               |
            |     HTTP metrics or               |
            |     ICMP metrics or               |
            | ------ | ------ | ------ | ------ | 

        - Metric:
            |  byte4 |  byte3 |  byte2 |  byte1 | 
            | ------ | ------ | ------ | ------ |
            |          pkt_rx                   |
            |                                   |
            | ------ | ------ | ------ | ------ | 
            |          pkt_tx                   |
            |                                   |
            | ------ | ------ | ------ | ------ | 
            |          byte_rx                  |
            |                                   |
            | ------ | ------ | ------ | ------ | 
            |          byte_tx                  |
            |                                   |
            | ------ | ------ | ------ | ------ | 
            |          bps_rx                   |
            |                                   |
            | ------ | ------ | ------ | ------ | 
            |          bps_tx                   |
            |                                   |
            | ------ | ------ | ------ | ------ | 
            |          err_rx                   |
            |                                   |
            | ------ | ------ | ------ | ------ | 
            |          err_tx                   |
            |                                   |
            | ------ | ------ | ------ | ------ | 
           
### 7. error response messages
    1. Error Response
        |  byte4 |  byte3 |  byte2 |  byte1 | 
        | ------ | ------ | ------ | ------ |
        | command         | length          |
        | ------ | ------ | ------ | ------ |\
        | char4  |  char3 |  char2 |  char1 | |
        | ------ | ------ | ------ | ------ | |  error
        | charN  |      .....               | | 
        | ------ | ------ | ------ | ------ |/


    2. Message Error Response
        |  byte4 |  byte3 |  byte2 |  byte1 | 
        | ------ | ------ | ------ | ------ |
        | command         | length          |
        | ------ | ------ | ------ | ------ |
        |             flow id               |
        | ------ | ------ | ------ | ------ |\
        | char4  |  char3 |  char2 |  char1 | |
        | ------ | ------ | ------ | ------ | |  error
        | charN  |      .....               | | 
        | ------ | ------ | ------ | ------ |/
