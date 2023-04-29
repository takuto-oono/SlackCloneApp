import classes from '@styles/Header.module.css'
import { Logout } from "@components/main/user";
import { useCookies } from 'react-cookie';

const Header = () => {
  const [cookies, setCookie, removeCookie] = useCookies(['token', 'user_id']);
  if (cookies.token) {
    // console.log(cookies.token)
    return (
      <header className={classes.header}>
        <h1>header</h1>
        <Logout />
      </header>
    );
  } else {
    return (
      <header className={classes.header}>
        <h1>header</h1>
      </header>
    );
  }
};

export default Header;
