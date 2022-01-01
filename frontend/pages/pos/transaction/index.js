import { useContext,useState,useEffect } from "react";
import DefaultLayout  from "@/layouts/default";
import axiosClient from "@/librarys/axiosClient"
import { ToastError,ToastSuccess } from "@/librarys/toaster";
import { Formik,Form,Field, ErrorMessage,useFormik,FieldArray } from 'formik';
import * as Yup from 'yup';

export default function Transaction() {
  let [customers,setCustmers] = useState([])
  let [products,setProducts] = useState([]);
  let [total,setTotal] = useState(0);

  const { values, errors, handleChange, setFieldValue,setValues } = useFormik({
    initialValues: {
     customer_id : 0,
     detail_orders : []
    }
  });

  const onAdd = () => {
    setValues({
      ...values,
      detail_orders : [...values.detail_orders,{product_id : 0,stock : 0,qty : 0,price : 0,total : 0}]
    })
  }

  const onDelete = (index) => {
    setValues({
      ...values,
      detail_orders : values.detail_orders.filter((_,indexItem) => indexItem != index)
    });
  }

  const onChooseProduct = (evt,indexItem) => {
    const detail_orders = values.detail_orders.map((item,index) => {
      if(indexItem == index){
        let {stock,price} = products.find(itemProduct => itemProduct.id == parseInt(evt.target.value))

        return {
          ...item,stock,price,
          product_id: parseInt(evt.target.value),          
          total : item.qty > 0 ? (item.qty*price) : 0
        }
      }else{
        return item;
      }
    })

    setValues({
      ...values,
      detail_orders : detail_orders
    })  
  }
  
  const onChangeQty = (evt,indexItem) => {    
    const detail_orders = values.detail_orders.map((item,index) => {
      if(indexItem == index){      
        return {
          ...item,
          qty : parseInt(evt.target.value),
          total : item.product_id > 0 ? (parseInt(evt.target.value)*item.price) : 0
        }
      }else{
        return item;
      }
    })

    setValues({
      ...values,
      detail_orders : detail_orders
    })
  }


  const onSubmit = async (values,{setErrors,setSubmitting}) => {      
    try{     
      await axiosClient.post("/transaction",{
        ...values,
        customer_id : parseInt(values.customer_id)
      });      

      ToastSuccess(`Berhasil Tambah Tranasaksi`)        
    }catch(err){                 
      console.log(err)
      ToastError(err)                              
    }finally{
      setSubmitting(false)
    }  
  }

  useEffect(() => {
    (async ()=> {
      try{
        let responseCustomer = await axiosClient.get("/customer/all");
        setCustmers(responseCustomer.data.data);

        let responseProduct = await axiosClient.get("/product/all");
        setProducts(responseProduct.data.data)
      }catch(err){
        console.log(err)
        ToastError("Terjadi Kesalahan")
      }
    })()
  }, [])

  useEffect(() => { 
    setTotal(values.detail_orders.reduce((amount,item) => amount+item.total,0));
  },[values.detail_orders])

  return (
    <DefaultLayout>
      <h1>Transaction</h1>

      <Formik
        initialValues={values}
        enableReinitialize={true}
        validationSchema={() => Yup.lazy((values) => {          
          return Yup.object()
          .shape({            
            // name : Yup.string().required('Required'),        
            // email : Yup.string().required("Required"),
            // password : !values.isEditable ? Yup.string().required("Required") : Yup.string(),                
          });    
        })}
        onSubmit={onSubmit}>     
          {({isSubmitting}) => (          

      <Form>

      <div>        
        <label>Customer</label>
        <Field name="customer_id" as="select">
            <option value="0">Pilih</option>
            {customers.map(item => (
              <option key={item.id} value={item.id}>{item.name}</option>                
            ))}
        </Field>
      
        <a href="#" onClick={onAdd}>Add</a>

        <table>
          <tr>
            <td>Product</td>
            <td>Stock</td>
            <td>Qty</td>
            <td>Price</td>
            <td>Total</td>
            <td>Opsi</td>
          </tr>

          <FieldArray
            name="detail_orders"
            render={() => (              
              <tbody>
              {values.detail_orders.map((detail_order, index) => (
              <tr key={index}>
                <td>
                  <Field name={`detail_orders[${index}].product_id`} as="select"
                    onChange={(event) => onChooseProduct(event,index)}>
                    <option value="0" disabled="true">Pilih</option>
                      {products.map(item => (
                        <option key={item.id} value={item.id}>{item.name}</option>                
                    ))}
                  </Field>
                </td>       
                <td>                 
                  <Field name={`detail_orders[${index}].stock`} readonly="true"/>
                </td> 
                <td>
                  <Field name={`detail_orders[${index}].qty`} 
                    onChange={(event) => onChangeQty(event,index)}/>
                </td>
                <td>
                  <Field name={`detail_orders[${index}].price`} readonly="true"/>
                </td>
                <td>
                  <Field name={`detail_orders[${index}].total`} readonly="true"/>
                </td>
                <td>
                  <a href="#" onClick={() => onDelete(index)}>Delete</a>
                </td>
              </tr>
              ))}   
              </tbody>
          )}/>

          <tr>
            <td colSpan="5"></td>                       
            <td>Total : {total}</td>
          </tr>
        </table>

        <div>
          <button 
            type="submit" 
            disabled={isSubmitting}>
              {isSubmitting ? '...' : 'Submit'}
          </button>    
        </div>
      </div>
      </Form>
        )}
      </Formik>
    </DefaultLayout>
  )
}
