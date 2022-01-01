import { useContext,useEffect } from "react";
import AppContext from "@/contexts/state"
import Link from "next/link";

export default function Home() {
  const auth = useContext(AppContext);

  useEffect(() => {
    console.log(auth);
  }, [])

  return (
    <>
     Hello Guest
     {auth.auth && <ul>
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
    </ul> }
    </>
  )
}
