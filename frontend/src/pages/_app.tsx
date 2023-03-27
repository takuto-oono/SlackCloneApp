import '@/styles/globals.css'
import { Routes, Route } from "react-router-dom";
import Login from '@src/pages/login_form';
import SignUp from '@src/pages/signUp_form';
import { Home } from '@src/pages/Home';
import IndexW from '@src/pages/workspace/workspace';
import CreateW from '@src/pages/workspace/create';
import ShowW from '@src/pages/workspace/show/[id]';
import Header from '@src/pages/header';

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
