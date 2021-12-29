import { ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import { AppWrapper } from '@/contexts/state'; 
import { getCookies,setCookies,removeCookies } from 'cookies-next';
import axiosServer from "@/librarys/axiosServer";
import axiosClient from "@/librarys/axiosClient";

function MyApp({ Component, pageProps,user}) {
  return (
    <AppWrapper user={user}>  
      <Component {...pageProps} />
      <ToastContainer/>
    </AppWrapper>
  )
}

MyApp.getInitialProps = async ({ctx}) => {
  let user = null;

  if(typeof window === 'undefined' && getCookies(ctx,"token").token){
    let response = await axiosServer(ctx).get("/me");
    user = response.data
  }

  if(typeof window !== 'undefined' && getCookies(null,'token').token){  
    let response = await axiosClient.get("/me");
    user = response.data     
  }

  return {
    user
  }
}

export default MyApp
