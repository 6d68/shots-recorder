import React, {ReactNode} from 'react'
import Head from 'next/head'
import Nav from "./Nav";


type Props = {
    children?: ReactNode
    title?: string
}

const Layout = ({children, title = 'Shots Recorder'}: Props) => (
    <div>
        <Head>
            <title>{title}</title>
            <meta charSet="utf-8"/>
            <meta name="viewport" content="initial-scale=1.0, width=device-width"/>
        </Head>
        <Nav/>
        <div className="bg-white">
            <div className="container mx-auto px-4">
                {children}
            </div>

        </div>
    </div>
)

export default Layout