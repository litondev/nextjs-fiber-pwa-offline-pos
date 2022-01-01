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
    var {data}  = await axiosServer(context).get("/order",{params});        
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
  const [orders,setOrders] = useState({...props.data})
  const [params,setParams] = useState({...props.params})
  const [show,setShow] = useState({
    isShow : false
  })

  useEffect(() => {
    setOrders({
      ...props.data
    })
  },[props.data])


  const onPage = (isNext) => {  
    router.push({
      pathname: router.pathname,
      query: { 
        ...props.params,
        page: isNext ? parseInt(orders.page) + 1 : parseInt(orders.page) - 1
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

  const onShow = (item) => {
    console.log(item)
    setShow({...item,isShow : true})
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
      <h1>Order</h1>

      {show.isShow && <div>
        <table>
          <tr>
            <td>Nama Pelanggan</td>
            <td>{show.customer.name}</td>
          </tr>
          <tr>
            <td>User</td>
            <td>{show.user.name}</td>
          </tr>
          <tr>
            <td>Total</td>
            <td>{show.total}</td>
          </tr>        
        </table>

        <h2>Detail</h2>

        {show.detail_orders.length && <table>
          <tr>
            <td>Nama Product</td>
            <td>Qty</td>
            <td>Price</td>            
            <td>Total</td>
          </tr>

          {show.detail_orders.map((item) => (
            <tr>
              <td>{item.product}</td>
              <td>{item.qty}</td>
              <td>{item.price}</td>
              <td>{parseInt(item.qty) * parseInt(item.price)}</td>
            </tr>
          ))}

          <tr>
            <td colspan="3"></td>
            <td>
            Total : {show.detail_orders.reduce((amount,item) => {
              return parseInt(item.qty)*parseInt(item.price)+amount
            },0)}
            </td>
          </tr>
        </table>
        }
      </div>
      }

      <div>
        <div>
          Searching 
          <input type="text" onKeyUp={onSearch}/>
        </div>

        <table>
          <thead>
          <tr>
            <td>Kategori</td>
            <td>Pelanggan</td>
            <td>Total</td>
            <td>Opsi</td>
          </tr>
          </thead>
          <tbody>
          {orders.data.map(item => (
            <tr key={item.id}>
              <td>{item.user.name}</td>
              <td>{item.customer.name}</td>
              <td>{item.total}</td>
              <td>
                <button onClick={() => onShow(item)}>Show</button>
              </td>
            </tr>
          ))}
          {!orders.data.length && (
            <tr>
              <td colSpan="3">
                Data tidak ditemukan
              </td>
            </tr>
          )}
          </tbody>
        </table>

        {orders.page > 1 && <button onClick={() => onPage(0)}>
          Sebelumnya
        </button> }

        {orders.total_page > orders.page && <button onClick={() => onPage(1)}>
          Selanjutnya
        </button>}
      </div>     
    </DefaultLayout>
  )
}
