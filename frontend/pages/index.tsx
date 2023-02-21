// import TestAPI1 from "./api/test_api_1";
import { User, getUsers, getUserById, postUser, updateUser, deleteUser, } from "./fetchAPI/user";
import LoginForm from "./component/login_form";
import { currentUser, login, } from "./fetchAPI/login";
import testAPI1 from "./api/test_api_1";
export default function Home() {

  // console.log(getUsers())
  let user: User = {
    id: "",
    name: "test docker",
    password: "test docker"
  }
  // let userRes = postUser(user)
  // console.log(userRes)
  return (
    <main>
      <h1>hello nextjs</h1>
      <h2>login</h2>
      < LoginForm />
    </main>
  )
}
