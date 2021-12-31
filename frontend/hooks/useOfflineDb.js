import { useState ,useEffect} from "react";

const useOfflineDb = (version = 1) => {
  const [db,setDb] = useState(null)

  const onError = (event)  => {
    console.log(event.target.error)
  }

  const onSuccess = (event) => {
    let database = event.target.result;

    database.onerror = (event) => onError(event);

    setDb(database)
  }

  // const onUpgradeNeeded = (event) => {    
  //   let database = event.target.result;

  //   database.onerror = (event) => onError(event)

  //   console.log("Fungsi ini akan terpangil jika ada perubahan versi db");

  //   try{
  //       let objectStoreCustomers = database.createObjectStore("categories", { keyPath: "ssn" , autoIncrement : true });
  //       objectStoreCustomers.createIndex("id", "id", { unique: false });
  //       objectStoreCustomers.createIndex("name", "name", { unique: false });
  //       objectStoreCustomers.createIndex("description", "description", { unique: false });
  //   }catch(err){
  //       console.log(err)
  //   }    
  // }

  useEffect(() => {    
    if (!window.indexedDB) {
        console.log("Your browser doesn't support a stable version of IndexedDB. Such and such feature will not be available.");
    }else{
        let requestDb = window.indexedDB.open("pos",4)
        requestDb.onerror = (event) => onError(event)
        requestDb.onsuccess = (event) => onSuccess(event)
        
        return () => {
          // requestDb.result.close();
        }
    }
  }, [])

  return [db];
};

export default useOfflineDb;
