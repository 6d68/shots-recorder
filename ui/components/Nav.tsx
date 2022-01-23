import {useUser} from "@auth0/nextjs-auth0";

const Nav = () => {

    const {user} = useUser();

    return (
        <div className="bg-white">
            <div className="container mx-auto px-4">
                <div className="flex items-center justify-between py-4">
                    <div>Logo</div>
                    {user ? (
                        <a href="/api/auth/logout"
                           className="text-gray-800 text-sm font-semibold border px-4 py-2 rounded-lg hover:text-orange-600 hover:border-orange-600">Logout</a>
                    ) : <a href="/api/auth/login"
                           className="text-gray-800 text-sm font-semibold border px-4 py-2 rounded-lg hover:text-orange-600 hover:border-orange-600">Login</a>
                    }
                </div>
            </div>
        </div>
    )
}

export default Nav