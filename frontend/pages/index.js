import { useContext,useEffect } from "react";
import AppContext from "@/contexts/state"

export default function Home() {
  const auth = useContext(AppContext);

  useEffect(() => {
    console.log(auth);
  }, [])

  return <h1>Hello World</h1>
}
