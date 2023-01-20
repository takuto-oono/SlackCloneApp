// import TestAPI1 from "./api/test_api_1";
import { getUsers, getUserById, postUser, updateUser, deleteUser, testAPI } from "./fetchAPI/user";
import testAPI1 from "./api/test_api_1";
export default function Home() {

  // console.log(getUsers())
  // TestAPI1()
  return (
    <main>
      <h1>hello nextjs</h1>
      <button onClick={testAPI}>テスト</button>
    </main>
  )
}