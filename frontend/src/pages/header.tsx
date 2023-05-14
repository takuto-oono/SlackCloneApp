import classes from '@styles/Header.module.css'
import { Logout } from "@components/main/user";
import { useCookies } from 'react-cookie';

const Header = () => {
  const [cookies, setCookie, removeCookie] = useCookies(['token', 'user_id']);
  if (cookies.token) {
    // console.log(cookies.token)
    return (
      <header className={classes.header}>
        <div className={classes.container}>
          <div className={classes.item}>
            <h3 className='text-pink-700'>header</h3>
          </div>
          <div className={classes.item}>
            <Logout />
          </div>
        </div>
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
