import axios from "axios";
import { getCookies,setCookies,removeCookies } from 'cookies-next';

export default function axiosServer(context){
	var instance = axios.create({
    	baseURL : process.env.NEXT_PUBLIC_API_URL,
    	headers : {
      		'Authorization' : getCookies(context,'token').token ? 'Bearer '+getCookies(context,'token').token : null
    	}
  	})

	function refreshToken(){	
		console.log("Refresh Token");
	}

	instance.interceptors.response.use(res => {	
		if(!res.data.access_token){	
			refreshToken();
		}

		return res;
	});

	return instance;
}