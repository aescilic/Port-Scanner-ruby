package main

import (
	"fmt"
	"net"
	"sort"
	"sync"
	"time"
)
const banner = `
  ____         ____                        
 / ___|  ___  / ___|  ___ __ _ _ __        
| |  _  / _ \ \___ \ / __/ _` | '_ \       
| |_| || (_) | ___) | (_| (_| | | | |      
 \____| \___/ |____/ \___\__,_|_| |_|      
                                           
 >> Go Port Scanner v1.0 | @aesclic
-------------------------------------------`

func main() {

	fmt.Println(banner)
	fmt.Printf("[*] Tarama Başlatıldı: %s\n", time.Now().Format("15:04:05"))
	fmt.Println("-------------------------------------------")

	target := "scanme.nmap.org" 
	fmt.Printf("[i] Hedef: %s\n", target)
	fmt.Println("[i] İlk 1024 port taranıyor, worker'lar hazırlanıyor...")

	
	var openPorts []int
	var wg sync.WaitGroup 

	
	portsCh := make(chan int, 100) 
	resultsCh := make(chan int)    


	for i := 0; i < cap(portsCh); i++ {
		go func() {
        for p := range portsCh {
				address := fmt.Sprintf("%s:%d", target, p)
				
				
				conn, err := net.DialTimeout("tcp", address, 2*time.Second)
				
				if err != nil {
					
					resultsCh <- 0 
					wg.Done() 
					continue
				}
				
				
				conn.Close()
				resultsCh <- p 
				wg.Done() 
			}
		}()
	}


	go func() {
		for r := range resultsCh {
			if r != 0 {
				openPorts = append(openPorts, r)
			}
			
		}
	}()
    for i := 1; i <= 1024; i++ {
		wg.Add(1) 
		portsCh <- i
	}

	
	wg.Wait()
	
	
	close(portsCh)
	close(resultsCh)

	
	fmt.Println("\n-------------------------------------------")
	fmt.Printf("[*] Tarama Tamamlandı: %s\n", time.Now().Format("15:04:05"))
	
	if len(openPorts) == 0 {
		fmt.Println("[-] Kötü haber: Hiçbir açık port bulunamadı.")
	} else {
		sort.Ints(openPorts)
		fmt.Printf("[+] Harika! Toplam %d açık port bulundu:\n", len(openPorts))
		for _, port := range openPorts {
			fmt.Printf("    [->] %d/tcp AÇIK\n", port)
		}
	}
	fmt.Println("-------------------------------------------")
	fmt.Println("Güle güle kullan! :)")
}
