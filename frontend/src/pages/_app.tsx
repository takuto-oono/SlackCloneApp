import '@styles/globals.css'
import { Routes, Route } from "react-router-dom";
import Login from '@pages/login_form';
import SignUp from '@pages/signUp_form';
import  Home  from '@pages/Home';
import IndexW from '@pages/workspace/workspace';
import CreateW from '@pages/workspace/create';
import ShowW from '@pages/workspace/show/[id]';
import Header from '@pages/header';

function App() {
  return (
    <div className="App">
      <Header />
      <Routes >
        <Route path="/">
          <Route index element={<Home />} />
          <Route path="login_form" element={<Login />} />
          <Route path="signUp_form" element={<SignUp />} />
          <Route path="workspace" >
            <Route index element={<IndexW />} />
            <Route path="create" element={<CreateW />} />
            <Route path={"show/:id"} element={<ShowW />}/>
          </Route>
        </Route>
      </Routes>
    </div>
  )
}

export default App;
