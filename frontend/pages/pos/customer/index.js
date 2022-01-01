import { useEffect,useState} from "react";
import DefaultLayout  from "@/layouts/default";
import axiosServer from "@/librarys/axiosServer";
import axiosClient  from "@/librarys/axiosClient";
import { useRouter } from 'next/router'
import {ToastError,ToastSuccess} from "@/librarys/toaster"
import { Formik,Form,Field, ErrorMessage,useFormik } from 'formik';
import * as Yup from 'yup';

const SigninSchema = Yup.object()
    .shape({            
        name : Yup.string().required('Required'),
        address : Yup.string().required('Required'),
        email : Yup.string().required("Required"),
        phone : Yup. string().required("Required"),                
    });
    
export async function getServerSideProps(context) {
  let isSuccess = true 

  let params = {
    search : context.query.search || '',
    per_page : context.query.per_page || 10,
    page : context.query.page || 1
  }
  
  try{
    var {data}  = await axiosServer(context).get("/customer",{params});        
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
  const [customers,setCustomers] = useState({...props.data})
  const [params,setParams] = useState({...props.params})
  const [loadings,setLoadings] = useState({
    isDelete : false,
  })


  useEffect(() => {
    console.log(props.data)
    setCustomers({
      ...props.data
    })
  },[props.data])


  const { values, errors, handleChange, setFieldValue,setValues } = useFormik({
    initialValues: {
      id : 0,
      isEditable : false,
      name : '', 
      address : '',
      email : '',
      phone : ''
    }
  });

  const onPage = (isNext) => {  
    router.push({
      pathname: router.pathname,
      query: { 
        ...props.params,
        page: isNext ? parseInt(customers.page) + 1 : parseInt(customers.page) - 1
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

  const onDelete = async (item) => {
    if(loadings.isDelete) return;

    try{
      setLoadings({
        ...loadings,
        isDelete : true
      })

      await axiosClient.delete("/customer/"+item.id)

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
        address : '',
        phone : '',
        email : '',
        isEditable : false
    })    
  }

  const onSubmit = async (values,{setErrors,setSubmitting}) => {      
    try{
      if(values.isEditable){
        await axiosClient.put("/customer/"+values.id,values)
      }else{
        await axiosClient.post("/customer",values);
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
      <h1>Customer</h1>

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

              <label>Email</label>
              <Field type="text" name="email"/>
              <ErrorMessage  
                name="email" 
                component="div" />
              <br/>

              <label>Phone</label>
              <Field type="text" name="phone"/>
              <ErrorMessage  
                name="phone" 
                component="div" />
              <br/>

              <label>Alamat</label>
              <Field type="textarea" name="address" as="textarea"></Field>
              <ErrorMessage  
                name="address" 
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
            <td>Phone</td>
            <td>Alamat</td>
            <td>Opsi</td>
          </tr>
          </thead>
          <tbody>
          {customers.data.map(item => (
            <tr key={item.id}>
              <td>{item.name}</td>
              <td>{item.email}</td>
              <td>{item.phone ?? '-'}</td>
              <td>{item.address ?? '-'}</td>
              <td>
                <button onClick={() => onEdit(item)}>Edit</button>
                <button onClick={() => onDelete(item)}>
                  {loadings.isDelete ? '...' : 'Delete'}
                </button>
              </td>
            </tr>
          ))}
          {!customers.data.length && (
            <tr>
              <td colSpan="3">
                Data tidak ditemukan
              </td>
            </tr>
          )}
          </tbody>
        </table>

        {customers.page > 1 && <button onClick={() => onPage(0)}>
          Sebelumnya
        </button> }

        {customers.total_page > customers.page && <button onClick={() => onPage(1)}>
          Selanjutnya
        </button>}
      </div>     
    </DefaultLayout>
  )
}
