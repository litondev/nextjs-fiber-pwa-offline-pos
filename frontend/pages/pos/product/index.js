import { useContext,useEffect } from "react";
import AppContext from "@/contexts/state"
import DefaultLayout  from "@/layouts/default";

export default function Home() {
  const auth = useContext(AppContext);

  useEffect(() => {
    console.log(auth);
  }, [])

  return (
    <DefaultLayout>
      <h1>Product</h1>
    </DefaultLayout>
  )
}
