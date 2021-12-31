import { ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import { AppWrapper } from '@/contexts/state'; 
import { getCookies,setCookies,removeCookies } from 'cookies-next';
import axiosServer from "@/librarys/axiosServer";
import axiosClient from "@/librarys/axiosClient";
import { Router } from 'next/router';
import IndexedDBProvider, { useIndexedDBStore } from "use-indexeddb";

const idbConfig = {
    databaseName: "pos",
    version: 4,
    stores: [
      {
        name: "customers",
        id: { keyPath: "ssn", autoIncrement: true },
        indices: [
          { name: "id", keyPath: "id", options: { unique: false } },
          { name: "name", keyPath: "name", options: { unique: false } },
          { name: "description", keyPath: "description" ,options: { unique: false } },
        ],
      },
    ],
  };

function MyApp({ Component, pageProps,user}) {
  return (
    <IndexedDBProvider config={idbConfig}>            
      <AppWrapper user={user}>  
        <Component {...pageProps} />
        <ToastContainer/>
      </AppWrapper>
    </IndexedDBProvider>    
  )
}

MyApp.getInitialProps = async ({ctx}) => {
  let user = null;

  if(typeof window === 'undefined' && getCookies(ctx,"token").token){
    try{
      let response = await axiosServer(ctx).get("/me");
      user = response.data
    }catch(err){
      if(err.response && err.response.status === 401){
        removeCookies("token",ctx)	
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
      if(err.response && err.response.status === 401){
        removeCookies("token")
        Router.push("/auth/login")
      }
    }
  }

  return {
    user
  }
}

export default MyApp
