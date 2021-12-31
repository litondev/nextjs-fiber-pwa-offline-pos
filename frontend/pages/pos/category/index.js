import { useEffect,useState} from "react";
import DefaultLayout  from "@/layouts/default";
import axiosServer from "@/librarys/axiosServer";
import axiosClient  from "@/librarys/axiosClient";
import { useRouter } from 'next/router'
import {ToastError,ToastSuccess} from "@/librarys/toaster"
import useOfflineDb from "@/hooks/useOfflineDb"
import { useIndexedDBStore } from "use-indexeddb";
import { Formik,Form,Field, ErrorMessage,useFormik } from 'formik';
import * as Yup from 'yup';

const SigninSchema = Yup.object()
    .shape({            
        name : Yup.string().required('Required'),
        description : Yup.string().required('Required'),
    });
    
export async function getServerSideProps(context) {
  let isSuccess = true 
  let isOffline = false;

  let params = {
    search : context.query.search || '',
    per_page : context.query.per_page || 10,
    page : context.query.page || 1
  }
  
  try{
    var {data}  = await axiosServer(context).get("/category",{params});        
  }catch(err){     
    isSuccess = false    

    if(err.response){
      let { status, statusText, data: response } = err.response
      var data = { status,statusText,data : response }
    }

    if(err.isAxiosError && !err.response){
      var data = { data : []}
      isSuccess = true;
      isOffline = true;
    }
  }
  
  return {
    props: {
      isSuccess,
      data,
      params,
      isOffline
    }, 
  }
}


export default function Category(props) {
  const router = useRouter();
  const [categories,setCategories] = useState({...props.data})
  const [params,setParams] = useState({...props.params})
  const [isOffline,setMode] = useState(props.isOffline ? true : false)
  const [loadings,setLoadings] = useState({
    isDelete : false,
    isForm : false,
    isPage : false,
  })

  const { add,openCursor,deleteAll ,update,deleteByID} = useIndexedDBStore("customers");  
  const [db] = useOfflineDb(1)
  const { values, errors, handleChange, setFieldValue,setValues } = useFormik({
    initialValues: {
      id : 0,
      isEditable : false,
      name : '', 
      description : ''
    }
  });

  useEffect(() => {
    if(isOffline && db){
      onLoadOffline()
    }
  },[db])

  useEffect(() => {
    setCategories({
      ...props.data
    })
  },[props.data])

  const onLoadOffline = (limit = 5,page = 1) => {        
    // onResetForm()
    let result = []
    let count = 0
    let getSerach = params.search ? 0 : limit;

    if(db){
      console.log(db);
      const transaction = db.transaction(['customers'], 'readonly');
      const objectStore = transaction.objectStore('customers');

      var countRequest = objectStore.count();
      countRequest.onsuccess = function() {
        let total_page = Math.ceil(countRequest.result / limit)
        let start_page = (page - 1) * limit;
        console.log("Page : "+page)    
        console.log("Start Page : "+start_page)
        console.log("Total Page : "+total_page)  
        console.log("Limit : "+(limit+start_page))

        objectStore.openCursor(null,'prev').onsuccess = function(event) {          
          const cursor = event.target.result;
            
            if(cursor && (count <= (limit+start_page) || getSerach < limit)) {                                            
              if(count >= start_page || getSerach < limit){
                if(cursor.value.name.indexOf(params.search) >= 0){                  
                  if(count >= start_page){
                    result.push(cursor.value)      
                    getSerach++;              
                  }
                }
              }

              cursor.continue();
              count++
            } else {                        

              setCategories({
                data : result,
                page : page,
                total_page : total_page,
                per_page  : limit
              })
              console.log('Entries displayed backwards.');
            }
        }
      }
    }
  }


  const onNextOffline = () => {
    let page = categories.page + 1;

    if(page > categories.total_page){
      page = categories.total_page
    }

    onLoadOffline(5,page)
  }

  const onBeforeOffline = () => {
    let page = categories.page - 1;
    if(page < 1){
      page = 1;
    }

    onLoadOffline(5,page)
  }

  // 
  const onSearchOffline = (event) => {
    console.log(event.target.value)
    setParams({
      ...params,
      search : event.target.value
    })

    if(event.key == 'Enter'){      
      onLoadOffline(5,1)      
    }
  }  

  const onPage = (isNext) => {  
    router.push({
      pathname: router.pathname,
      query: { 
        ...props.params,
        page: isNext ? parseInt(categories.page) + 1 : parseInt(categories.page) - 1
      },
    });    
  }

  const onSearch = (event) => {
    console.log(event.target.value)
    setParams({
      ...params,
      search : event.target.value
    })

    if(event.key == 'Enter'){
      router.push({
        pathname: router.pathname,
        query: {
          ...params,
          page: 1
        }
      });          
    }
  }
  const onDeleteOffline = async (item) => {
    try{
      await deleteByID(item.ssn)
     
      ToastSuccess("Berhasil menghapus data")

      onLoadOffline()
    }catch(err){
      console.log(err)
      ToastError(err)
    }      
  }

  const onDelete = async (item) => {
    if(loadings.isDelete) return;

    try{
      setLoadings({
        ...loadings,
        isDelete : true
      })
      await axiosClient.delete("/category/"+item.id)
      ToastSuccess("Berhasil menghapus data")
      setMode(false)
      router.push({
        pathname: router.pathname,
        query: props.params
      });        
    }catch(err){
      if(err.isAxiosError && !err.response){
        setMode(true)
        onDeleteOffline(item)        
      }else{
        console.log(false)
        ToastError(err)   
      }                     
    }finally{
      setLoadings({
        ...loadings,
        isDelete : false
      })
    }
  }
   
  const onResetForm = () => {
    setValues({
        id : 0,
        name : '',
        description : '',
        isEditable : false
    })    
  }

  const onSubmitOffline = async(values) => {
    try{
      if(!values.isEditable){
        await add({ 
          name:  values.name,
          description : values.description,                            
        })
      }else{
        await update({
          ssn : values.ssn,
          name : values.name,
          description : values.description
        })
      } 
     
      console.log(values)

      ToastSuccess(`Berhasil ${values.isEditable ? 'Edit' : 'Tambah'} Data`)

      onLoadOffline()
    }catch(err){
      console.log(err)
      ToastError(err)
    }      
  }

  const onSubmit = async (values,{setErrors,setSubmitting}) => {      
    try{
      if(values.isEditable){
        await axiosClient.put("/category/"+values.id,values)
      }else{
        await axiosClient.post("/category",values);
      }

      ToastSuccess(`Berhasil ${values.isEditable ? 'Edit' : 'Tambah'} Data`)

      setMode(false)

      onResetForm()

      router.push({
        pathname: router.pathname,
        query: params
      });        
    }catch(err){            
      if(err.isAxiosError && !err.response){
        setMode(true)
        onSubmitOffline(values)        
      }else{
        console.log(err)
        ToastError(err,setErrors)                        
      }
    }finally{
      setSubmitting(false)
    }  
  }
   
  const onEdit = (item) => {
    console.log(item)
    setValues({         
      ...item,
      isEditable : true
    })
  }

  if(!props.isSuccess){
    return (
      <DefaultLayout>
        <h1>Terjadi Kesalahan</h1>
      </DefaultLayout>
    )
  }

  return (
    <DefaultLayout>
      <h1>Category</h1>

      <div>
        Add/Edit 
        <Formik
          initialValues={values}
          enableReinitialize={true}
          validationSchema={SigninSchema}
          onSubmit={onSubmit}>        
          {({isSubmitting}) => (          
            <Form>
              <label>Nama</label>
              <Field  type="text" name="name"/>
              <ErrorMessage  
                name="name" 
                component="div" />
              <br/>
              <label>Deskripsi</label>
              <Field type="textarea" name="description" as="textarea"></Field>
              <ErrorMessage  
                name="description" 
                component="div" />
              <br/>
              <div>
                  <button 
                      type="submit" 
                      disabled={isSubmitting}>
                      {isSubmitting ? '...' : 'Submit'}
                  </button>
                  <button type="reset"
                      onClick={onResetForm}>
                      Reset
                  </button>
              </div>
            </Form>
          )}
        </Formik>
        <br/>
      </div>

      {!isOffline ? <div>
        <div>
          Searching 
          <input type="text" onKeyUp={onSearch}/>
        </div>

        <table>
          <thead>
          <tr>
            <td>Nama</td>
            <td>Deskripsi</td>
            <td>Opsi</td>
          </tr>
          </thead>
          <tbody>
          {categories.data.map(item => (
            <tr key={item.id}>
              <td>{item.name}</td>
              <td>{item.description}</td>
              <td>
                <button onClick={() => onEdit(item)}>Edit</button>
                <button onClick={() => onDelete(item)}>
                  {loadings.isDelete ? '...' : 'Delete'}
                </button>
              </td>
            </tr>
          ))}
          {!categories.data.length && (
            <tr>
              <td colSpan="3">
                Data tidak ditemukan
              </td>
            </tr>
          )}
          </tbody>
        </table>

        {categories.page > 1 && <button onClick={() => onPage(0)}>
          Sebelumnya
        </button> }

        {categories.total_page > categories.page && <button onClick={() => onPage(1)}>
          Selanjutnya
        </button>}
      </div>
      : <div>
      <div>
        Searching 
        <input type="text" onKeyUp={onSearchOffline}/>
      </div>

      <table>
        <thead>
        <tr>
          <td>Nama</td>
          <td>Deskripsi</td>
          <td>Opsi</td>
        </tr>
        </thead>
        <tbody>
        {categories.data.map(item => (
          <tr key={item.id}>
            <td>{item.name}</td>
            <td>{item.description}</td>
            <td>
              <button onClick={() => onEdit(item)}>Edit</button>
              <button onClick={() => onDelete(item)}>
                {loadings.isDelete ? '...' : 'Delete'}
              </button>
            </td>
          </tr>
        ))}
        {!categories.data.length && (
          <tr>
            <td colSpan="3">
              Data tidak ditemukan
            </td>
          </tr>
        )}
        </tbody>
      </table>

      {categories.page > 1 && <button onClick={() => onBeforeOffline()}>
          Sebelumnya
        </button> }

        {categories.total_page > categories.page && <button onClick={() => onNextOffline()}>
          Selanjutnya
        </button>}
    </div>}
    </DefaultLayout>
  )
}
