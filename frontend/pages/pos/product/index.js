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
    var {data}  = await axiosServer(context).get("/product",{params});        
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
  const [products,setProducts] = useState({...props.data})
  const [params,setParams] = useState({...props.params})
  const [loadings,setLoadings] = useState({
    isDelete : false,
  })
  
  const [showForm,setShowForm] = useState(false);
  const [categories,setCatgories] = useState([])

  useEffect(() => {
    (async () => {
      try{
        let response = await axiosClient.get("/category/all");           
        setCatgories(response.data.data)
      }catch(err){
        console.log(err)
        ToastError("Terjadi Kesalahan")
      }
    })()
  },[])

  useEffect(() => {
    setProducts({
      ...props.data
    })
  },[props.data])


  const { values, errors, handleChange, setFieldValue,setValues } = useFormik({
    initialValues: {
      id : 0,
      isEditable : false,      
      name : '',     
      code : '',
      description : '',
      stock : 0,
      price : 0,      
      photo : '',
      category_id : 0  
    }
  });

  const onPage = (isNext) => {  
    router.push({
      pathname: router.pathname,
      query: { 
        ...props.params,
        page: isNext ? parseInt(products.page) + 1 : parseInt(products.page) - 1
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

      await axiosClient.delete("/product/"+item.id)

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
        isEditable : false,
        name : '',     
        code : values.code,
        description : '',
        stock : 0,
        price : 0,
        photo : '',
        category_id : 0  
    })    
  }

  const onSubmit = async (values,{setErrors,setSubmitting}) => {      
    try{  
      let formData = new FormData(document.getElementById("form-product"))            

      if(values.isEditable){
        await axiosClient.put("/product/"+values.id,formData)
      }else{
        await axiosClient.post("/product",formData);
      }

      ToastSuccess(`Berhasil ${values.isEditable ? 'Edit' : 'Tambah'} Data`)

      setShowForm(false)  

      router.push({
        pathname: router.pathname,
        query: params          
      });        

      onResetForm()
    }catch(err){                 
      console.log(err)
      ToastError(err,setErrors)                              
    }finally{
      setSubmitting(false)
    }  
  }
   
  const onAdd = async () => {
    try{
      let response = await axiosClient.get("/product/code")

      setValues({
        ...values,
        code : response.data.code,
      });

      setShowForm(true)  
    }catch(err){
      console.log(err);
      ToastError(err)
    }
  }

  const onEdit = (item) => {
    setShowForm(true)
    setValues({         
      ...item,
      isEditable : true
    })
  }

  const onChangePhoto = (evt) => {
    if(!evt.target.files[0]){
        return false;
    } 

    if(!['image/jpeg','image/jpg','image/png'].includes(evt.target.files[0].type)){
        evt.target.value = "";              
        return false;
    }  
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
      <h1>Product</h1>

      {showForm && <div>
        Add/Edit 
        <Formik
          initialValues={values}
          enableReinitialize={true}        
          validationSchema={() => Yup.lazy((values) => {          
            return Yup.object()
            .shape({            
              code : Yup.string().required("Reuqired"),
              name : Yup.string().required('Required'),      
              price : Yup.number().required("Required"),
              category_id : Yup.number().required("Required")
            });    
          })}
          onSubmit={onSubmit}>        
          {({isSubmitting}) => (          
            <Form id="form-product">
              <label>Code</label>
              <Field  type="text" name="code" readonly="true"/>
              <br/>

              <label>Nama</label>
              <Field  type="text" name="name"/>
              <ErrorMessage  
                name="name" 
                component="div" />
              <br/>

              <label>Price</label>
              <Field type="text" name="price"/>
              <ErrorMessage  
                name="price" 
                component="div" />
              <br/>

              <label>Stock</label>
              <Field type="text" name="stock"/>
              <ErrorMessage  
                name="stock" 
                component="div" />
              <br/>          

              <label>Category</label>
              <Field name="category_id" as="select">
                  <option value="0">Pilih</option>
                  {categories.map(item => (
                    <option key={item.id} value={item.id} selected={item.id == values.category_id}>{item.name}</option>                
                  ))}
              </Field>
              
              <br/>

              <label>Deskripsi</label>
              <Field as="textarea" name="description"/>
              <ErrorMessage  
                name="description" 
                component="div" />
              <br/>

              <label>Photo</label>
              <input type="file" name="photo" onChange={onChangePhoto}/>

              <ErrorMessage  
                name="category_id" 
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
      </div> }         

      <div>
        {!showForm && <div>
            <button onClick={onAdd}>Tambah</button>
        </div>}

        <div>
          Searching 
          <input type="text" onKeyUp={onSearch}/>
        </div>

        <table>
          <thead>
          <tr>
            <td>Gambar</td>
            <td>Nama</td>
            <td>Harga</td>
            <td>Qty</td>
            <td>Kategori</td> 
            <td>Opsi</td>
          </tr>
          </thead>
          <tbody>
          {products.data.map(item => (
            <tr key={item.id}>
              <td>
                <img src={item.photo ? process.env.NEXT_PUBLIC_API_ASSET_URL + "/images/products/" + item.photo : process.env.NEXT_PUBLIC_API_ASSET_URL + "/images/products/default.png"}
                  width="100px"/></td>
              <td>{item.name}</td>
              <td>{item.price}</td>
              <td>{item.category.name}</td>
              <td>
                <button onClick={() => onEdit(item)}>Edit</button>
                <button onClick={() => onDelete(item)}>
                  {loadings.isDelete ? '...' : 'Delete'}
                </button>
              </td>
            </tr>
          ))}
          {!products.data.length && (
            <tr>
              <td colSpan="3">
                Data tidak ditemukan
              </td>
            </tr>
          )}
          </tbody>
        </table>

        {products.page > 1 && <button onClick={() => onPage(0)}>
          Sebelumnya
        </button> }

        {products.total_page > products.page && <button onClick={() => onPage(1)}>
          Selanjutnya
        </button>}
      </div>     
    </DefaultLayout>
  )
}
