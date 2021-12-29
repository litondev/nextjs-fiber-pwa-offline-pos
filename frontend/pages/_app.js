import { ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import { AppWrapper } from '@/contexts/state'; 
import { getCookies,setCookies,removeCookies } from 'cookies-next';
import axiosServer from "@/librarys/axiosServer";
import axiosClient from "@/librarys/axiosClient";
import { Router } from 'next/router';

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
    try{
      let response = await axiosServer(ctx).get("/me");
      user = response.data
    }catch(err){
      removeCookies("token",ctx)	
      if(err.response && err.response.status === 401){
        ctx.res.writeHead(302, {
          Location: '/auth/login'
        });      
        ctx.res.end();
      }
    }
  }

  if(typeof window !== 'undefined' && getCookies(null,'token').token){  
    try{
      let response = await axiosClient.get("/me");
      user = response.data     
    }catch(err){
      removeCookies("token")
      if(err.response && err.response.status === 401){
        Router.push("/auth/login")
      }
    }
  }

  return {
    user
  }
}

export default MyApp
