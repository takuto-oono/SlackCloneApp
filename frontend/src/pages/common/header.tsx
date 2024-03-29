import { Logout } from '@src/components/user/user'
import { ShowLoginUserName } from '@src/components/user/show_login_username'
import { loginUserState } from '@src/utils/atom'
import { useRecoilValue } from 'recoil'

const Header = () => {
  const loginUser = useRecoilValue(loginUserState)

  if (loginUser.length != 0) {
    return (
      <header className='bg-purple-200 px-2.5 py-2.5 border-b-2 border-pink-50'>
        <div className='h-full flex' id='container'>
          <div className='float-left px-8 py-5 text-center' id='item'>
            <button
              type='button'
              className='inline-block align-baseline font-bold text-pink-700  text-2xl hover:text-blue-800'
            >
              SlackCloneApp
            </button>
          </div>
          <div className='float-left px-8 py-5 text-center' id='item'>
            <Logout />
          </div>
          <div className='float-left px-8 py-5 text-center'>
            <ShowLoginUserName />
          </div>
        </div>
      </header>
    )
  } else {
    return (
      <header className='bg-purple-200 px-2.5 py-2.5 border-b-2 border-pink-50'>
        <div className='h-full flex' id='container'>
          <div className='float-left px-8 py-5 text-center' id='item'>
            <p className='text-pink-700 text-2xl'>SlackCloneApp</p>
          </div>
        </div>
      </header>
    )
  }
}

export default Header
