import classes from '@styles/Header.module.css'
import { Logout } from "@components/login";

const Header=()=> {
  return (
    <div className={classes.header}>
      <h1>header</h1>
      <Logout />
    </div>
	);
};

export default Header;
