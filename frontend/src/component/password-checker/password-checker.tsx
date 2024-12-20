import React from "react";
import WebService from "../../service/web/web-service";
import RouteEnum from "../../service/web/route-enum";

interface PasswordCheckerProps {
    webService: WebService;
}

const PasswordCheckerComponent: React.FC<PasswordCheckerProps> = ({ webService }) => {
    const [password, setPassword] = React.useState<string>("");
    const [isLeakedPassword, setIsLeakedPassword] = React.useState<boolean | null>(null);

    const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        const response = await webService.post(RouteEnum.CheckPassword, { password });
        const data = await response.json();
        setIsLeakedPassword(data.found);
    };

    return (
    <form onSubmit={handleSubmit} className="flex flex-col items-center bg-sky-200 shadow-md rounded px-8 pt-6 pb-8 mb-4">
        <h1 className="text-4xl font-bold text-white bg-sky-300 border-4 border-white w-full text-center mb-4 p-8 rounded">
        🎸 Have I Been Rocked? 🎸
        </h1>
        <div className="flex w-full">
            <input 
                type="text" 
                value={password} 
                onChange={(event) => setPassword(event.target.value)} 
                className="shadow appearance-none border rounded-l w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" 
                placeholder="Enter your password"
            />
            <button 
                type="submit" 
                className="bg-sky-400 hover:bg-sky-500 text-white font-bold py-2 px-4 rounded-r focus:outline-none focus:shadow-outline"
                disabled={password.length === 0}
            >
                Rocked?
            </button>
        </div>
          <div className="text-center mt-4 text-2xl">
            { isLeakedPassword !== null ? <>
              {!isLeakedPassword && <p className="text-green-500">Your password is not leaked ✔</p>}
              {isLeakedPassword && <p className="text-red-500">Your password is leaked ❌</p>}
              </> : <p className="text-blue-500">Enter your password to check if it is leaked 🔒</p>}
          </div>
    </form>
    );
};

export default PasswordCheckerComponent;