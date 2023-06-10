import { Logout } from "@components/main/user";
import { ShowLoginUserName } from "@src/components/header/show_login_username";
import { useCookies } from "react-cookie";

const Header = () => {
	const [cookies, setCookie, removeCookie] = useCookies(["token", "user_id"]);
	if (cookies.token) {
		return (
			<header className="bg-purple-200 px-2.5 py-2.5 border-b-2 border-pink-50">
				<div className="h-full flex" id="container">
					<div className="float-left px-8 py-5 text-center" id="item">
						<h3 className="text-pink-700  text-2xl">header</h3>
					</div>
					<div className="float-left px-8 py-5 text-center" id="item">
						<Logout />
					</div>
					<div className="float-left px-8 py-5 text-center">
						<ShowLoginUserName />
					</div>
				</div>
			</header>
		);
	} else {
		return (
			<header className="bg-purple-200 px-2.5 py-2.5 border-b-2 border-pink-50">
				<div className="h-full flex" id="container">
					<div className="float-left px-8 py-5 text-center" id="item">
						<p className="text-pink-700 text-2xl">header</p>
					</div>
				</div>
			</header>
		);
	}
};

export default Header;
