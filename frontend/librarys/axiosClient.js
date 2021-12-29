import axios from "axios";
import Router from "next/router";
import { getCookies,setCookies,removeCookies } from 'cookies-next';

axios.defaults.headers.common['X-Requested-With'] = 'XMLHttpRequest';
axios.defaults.baseURL = process.env.NEXT_PUBLIC_API_URL;

axios.interceptors.request.use(config => { 	
  	config.headers['Authorization'] = getCookies(null,'token') ? 'Bearer '+getCookies(null,'token') : null
  	return config;
});

function refreshToken(){	
	console.log("Refresh Token");
}

axios.interceptors.response.use(res => {	
	if(!res.data.access_token) refreshToken();	
	return res;
});

export default axios;