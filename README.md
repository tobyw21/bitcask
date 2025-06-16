BitCask
===========================

implementing according to the [paper](https://riak.com/assets/bitcask-intro.pdf)


idea...
```                                                                                          
  store := BitStore("store1")                                                                         
                                                                                                      
  find a valid Object ID                                                                              
          |                                                                                           
          |                                                                                           
          v                                                                                           
  Create store Oid                                                                                    
          |                                                                                           
          v               if KV entry not in catalog                                                  
  Create a KV enrty file  -------------------------> Create Oid, KV store name as a map {Oid: KV name}
          |                                            save the map into data/catalog file            
          |                                                           |                               
          v                                                           |                               
       read catalog data <--------------------------------------------+                               
          |                                                                                            
          v                                                                                            
   User Set(), Get(), Remove(), Update() KVdata                                                       
   Update catalog accordingly                                                                         
          |                                                                                           
          v                                                                                           
   Seek data pos in the data/Oid/KVData.oid                                                           
   (check crc??)                                                                                      
   write/append to the file                                                                           
   updata data crc                                                                                    
         |                                                                                            
         |                                                                                            
         |                                                                                            
         v                                                                                            
   close bitstore                                                                                     
         |                                                                                            
         |                                                                                            
         |                                                                                            
         v                                                                                            
   fsync all memory KeyDir data into file ------> Possible merge operation??                          
                                                  if merge update catalog                          
```

implemet a vfd pool? maybe useful?? in order to make less open() and close()

dir should always have ea `data/` which has `data/catalog` and `data/KVoid/KVdata.oid`, not sure if Oid system is a good idea...

no buffer pool, it's a lightweight KV store DB so no need.
no indexing, no need, search is (almost) linear O(1)