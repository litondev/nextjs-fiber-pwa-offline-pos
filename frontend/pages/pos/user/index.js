import { useEffect,useState} from "react";
import DefaultLayout  from "@/layouts/default";
import axiosServer from "@/librarys/axiosServer";
import axiosClient  from "@/librarys/axiosClient";
import { useRouter } from 'next/router'
import {ToastError,ToastSuccess} from "@/librarys/toaster"
import { Formik,Form,Field, ErrorMessage,useFormik } from 'formik';
import * as Yup from 'yup';


    
export async function getServerSideProps(context) {
  let isSuccess = true 

  let params = {
    search : context.query.search || '',
    per_page : context.query.per_page || 10,
    page : context.query.page || 1
  }
  
  try{
    var {data}  = await axiosServer(context).get("/user",{params});        
  }catch(err){     
    isSuccess = false    

    if(err.response){
      let { status, statusText, data: response } = err.response
      var data = { status,statusText,data : response }
    }
  }
  
  return {
    props: {
      isSuccess,
      data,
      params,  
    }, 
  }
}


export default function Category(props) {
  const router = useRouter();
  const [users,setUsers] = useState({...props.data})
  const [params,setParams] = useState({...props.params})
  const [loadings,setLoadings] = useState({
    isDelete : false,
  })


  useEffect(() => {
    console.log(props.data)
    setUsers({
      ...props.data
    })
  },[props.data])


  const { values, errors, handleChange, setFieldValue,setValues } = useFormik({
    initialValues: {
      id : 0,
      isEditable : false,
      name : '',     
      email : '',
      password : ''
    }
  });

  const onPage = (isNext) => {  
    router.push({
      pathname: router.pathname,
      query: { 
        ...props.params,
        page: isNext ? parseInt(users.page) + 1 : parseInt(users.page) - 1
      },
    });    
  }

  const onSearch = (event) => {
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

  const onDelete = async (item) => {
    if(loadings.isDelete) return;

    try{
      setLoadings({
        ...loadings,
        isDelete : true
      })

      await axiosClient.delete("/user/"+item.id)

      ToastSuccess("Berhasil menghapus data")

      router.push({
        pathname: router.pathname,
        query: props.params
      });        
    }catch(err){    
      console.log(false)
      ToastError(err)     
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
        email : '',
        password : '',
        isEditable : false
    })    
  }

  const onSubmit = async (values,{setErrors,setSubmitting}) => {      
    try{
      if(values.isEditable){
        await axiosClient.put("/user/"+values.id,values)
      }else{
        await axiosClient.post("/user",values);
      }

      ToastSuccess(`Berhasil ${values.isEditable ? 'Edit' : 'Tambah'} Data`)

      onResetForm()

      router.push({
        pathname: router.pathname,
        query: params          
      });        
    }catch(err){                 
      console.log(err)
      ToastError(err,setErrors)                              
    }finally{
      setSubmitting(false)
    }  
  }
   
  const onEdit = (item) => {
    console.log(item)
    setValues({         
      ...item,
      password : '',
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
      <h1>User</h1>

      <div>
        Add/Edit 
        <Formik
          initialValues={values}
          enableReinitialize={true}
        
          validationSchema={() => Yup.lazy((values) => {          
            return Yup.object()
            .shape({            
              name : Yup.string().required('Required'),        
              email : Yup.string().required("Required"),
              password : !values.isEditable ? Yup.string().required("Required") : Yup.string(),                
            });    
          })}
          onSubmit={onSubmit}>        
          {({isSubmitting}) => (          
            <Form>
              <label>Nama</label>
              <Field  type="text" name="name"/>
              <ErrorMessage  
                name="name" 
                component="div" />
              <br/>

              <label>Email</label>
              <Field type="text" name="email"/>
              <ErrorMessage  
                name="email" 
                component="div" />
              <br/>

              <label>Password</label>
              <Field type="text" name="password"/>
              <ErrorMessage  
                name="password" 
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

      <div>
        <div>
          Searching 
          <input type="text" onKeyUp={onSearch}/>
        </div>

        <table>
          <thead>
          <tr>
            <td>Nama</td>
            <td>Email</td>
            <td>Opsi</td>
          </tr>
          </thead>
          <tbody>
          {users.data.map(item => (
            <tr key={item.id}>
              <td>{item.name}</td>
              <td>{item.email}</td>
              <td>
                <button onClick={() => onEdit(item)}>Edit</button>
                <button onClick={() => onDelete(item)}>
                  {loadings.isDelete ? '...' : 'Delete'}
                </button>
              </td>
            </tr>
          ))}
          {!users.data.length && (
            <tr>
              <td colSpan="3">
                Data tidak ditemukan
              </td>
            </tr>
          )}
          </tbody>
        </table>

        {users.page > 1 && <button onClick={() => onPage(0)}>
          Sebelumnya
        </button> }

        {users.total_page > users.page && <button onClick={() => onPage(1)}>
          Selanjutnya
        </button>}
      </div>     
    </DefaultLayout>
  )
}
