package main


import (
  "fmt"
  "flag"
  "net/http"
  "time"
  "github.com/warthog618/gpio"
)

var host string
var poll_time string
var pinnum int
var reboot_wait string


func init (){
    
    flag.StringVar(&host, "host","", "Host ip to ping")
    flag.StringVar(&poll_time,"poll_time","30","Window to poll in seconds")
    flag.IntVar(&pinnum,"pin",4,"which pin to trigger relay")
    flag.StringVar(&reboot_wait,"reboot_wait","90","How long to wait for reboot after fail")
    
}



func main(){
    
    flag.Parse()
    
    go backgroundCheck()
    
    fmt.Println("Waiting for checks...")
    
    select{}
        
    
}

func backgroundCheck() {
	    
    dur,_ := time.ParseDuration(poll_time+"s")
    
    ticker := time.NewTicker(dur)
    
    err := gpio.Open()
    
    defer gpio.Close()
    
    if(err!=nil){}
    
    pin := gpio.NewPin(pinnum)
    pin.Output()  
    
	for _ = range ticker.C {
		
		    
		    
		    fmt.Printf("Checking %s...\n",host)
		    
		    status := urlChecker(host)
		    
            if(status == false){
                
                 //check again!
                 
                 time.Sleep(5 * time.Second)
                 
                 status2 := urlChecker(host)
                  
                  if(status2 == false ){ 
                  
                      fmt.Printf("Unable to reach host: %s\n",host)
                      pin.High()
                      time.Sleep(3 * time.Second)
                      pin.Low()
                      fmt.Printf("Waiting for reboot: %s\n",host)
                      
                      rb_w,_ := time.ParseDuration(reboot_wait+"s")
                      
                      time.Sleep(rb_w) //sleep for 90 seconds to allow reboot
                  
                  }
                     
            }else{
                
                fmt.Println("Check good!")
            }

		
	}
}


func urlChecker (h string) bool{ 
    
c := &http.Client{
    Timeout: 3 * time.Second,
}    
    
res, err := c.Get(h)

if err != nil {
  fmt.Println(err)
  return false   

}



if(res.StatusCode == 200){
    
    return true
    
}

  return false
    
    
}
