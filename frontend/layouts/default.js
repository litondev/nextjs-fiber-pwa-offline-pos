import Link from "next/link";
import { getCookies,setCookies,removeCookies } from 'cookies-next';
import Router from "next/router";

export default function DefaultLayout({children}){
    const onLogout = () => {
        removeCookies('token')
        Router.push('/')   
    }

    return (
        <>
            <ul>
                <li>
                    <Link href="/">Home</Link>                                                                
                </li>
                <li>
                    <Link href="/pos/category">Category</Link>                                                                
                </li>
                <li>
                    <Link href="/pos/product">Product</Link>
                </li>
                <li>
                    <Link href="/pos/customer">Customer</Link>
                </li>
                <li>
                    <Link href="/pos/order">Order</Link>
                </li>
                <li>
                    <Link href="/pos/transaction">Transaksi</Link>
                </li>
                <li>
                    <Link href="/pos/user">User</Link>
                </li>
                <li>
                    <a href="#"  onClick={onLogout}>
                        Keluar
                    </a>
                </li>
            </ul>

            <div>
                {children}
            </div>
        </>
    )
}