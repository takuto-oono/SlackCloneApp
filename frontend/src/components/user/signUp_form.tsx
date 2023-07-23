import Link from 'next/link'
import { useRouter } from 'next/router'
import React, { useState } from 'react'
import { signUp } from 'src/fetchAPI/user'

const SignUpForm = () => {
  const [name, setName] = useState('')
  const [password, setPassword] = useState('')
  const router = useRouter()

  const nameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setName(e.target.value)
  }

  const passwordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setPassword(e.target.value)
  }

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    console.log('signup')
    let user = { name: name, password: password }
    signUp(name, password)
  }

  return (
    <div>
      <form className='px-8 py-8' onSubmit={handleSubmit}>
        <p className='text-2xl p-1'>SingUp</p>
        <div className='mb-4'>
          <label className='block mb-2 font-bold'>名前</label>
          <input
            className='border border-black w-full py-2 px-3'
            type='text'
            value={name}
            name='name'
            onChange={nameChange}
            maxLength={80}
            required
          />
        </div>
        <div className='mb-6'>
          <label className='block mb-2 font-bold'>パスワード</label>
          <input
            className='border border-black w-full py-2 px-3'
            type='password'
            value={password}
            name='password'
            onChange={passwordChange}
            minLength={6}
            maxLength={72}
            required
          />
        </div>
        <div>
          <button
            className='bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline'
            type='submit'
          >
            作成
          </button>
        </div>
      </form>
      <div>
        <button
          type='button'
          onClick={() => router.push('/')}
          className='inline-block align-baseline font-bold text-sm text-blue-500 hover:text-blue-800'
        >
          既に作成してある方はこちらへ
        </button>
      </div>
    </div>
  )
}

export default SignUpForm
