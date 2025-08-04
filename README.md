BitCask
===========================

implementing according to this [paper](https://riak.com/assets/bitcask-intro.pdf)


ideas...
```                                                                                          
  store := BitStore("store1") --------------------- if store exists ---------------------------------------+                                                                     
           |                                                                                               |             
           |   if store does not exist                                                                     |
           v                                                                                               |
  find a valid Object ID                                                                                   |
          |                                                                                                |
          |                                                                                                |
          v                                                                                                |
  Create store Oid                                                                                         |
          |                                                                                                |
          v               if KV entry not in catalog                                                       |
  Create a KV enrty file  -------------------------> Create Oid, KV store name as a map {Oid: KV name}     |
          |                                            save the map into data/catalog file                 |
          |                                                           |                                    |
          v                                                           v                                    |
       read catalog data <---------------------------------------------------------------------------------+
          |                                                                                            
          v                                                                                            
   User Set(), Get(), Remove() KV data store                                                 
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

implemet a vfd pool? maybe useful?? in order to make less open() and close() or if we have multiple bitstore opened like

```go
        b1 := BitStore("store1")
        b2 := BitStore("store2")
        b3 := BitStore("store3")
        ...
```

should always have a `data/` dir which contains `data/catalog` and `data/KV.oid/KVdata.oid`, not sure if Oid system is a good idea...

or a store name based dir would be easier to implement

no buffer pool (for now), it's a lightweight KV store DB so no need.

no indexing, no need, search is (almost) linear O(1)

for merge, in the initial stage don't implement it, let's keep multiple data files in the data dir

TODO:

 - [ ] how to store `nextFree` in the hint file
 - [ ] read, write hint file
 - [ ] calculate crc
 - [ ] get function
 - [ ] delete function
 - [ ] merge (compact), later..
