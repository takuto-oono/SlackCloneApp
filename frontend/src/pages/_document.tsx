import { Html, Head, Main, NextScript } from 'next/document'

export default function Document() {
  return (
    <Html lang="en" className='h-full max-w-full max-h-screen overflow-x-hidden'>
      <Head />
      <body className='h-full max-w-full max-h-screen overflow-x-hidden'>
        <Main />
        <NextScript />
      </body>
    </Html>
  )
}
