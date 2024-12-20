import React, { useEffect } from "react";

interface PasswordGenerationConfig {
    length: number;
    uppercase: boolean;
    lowercase: boolean;
    numbers: boolean;
    symbols: boolean;
}

const DefaultConfig: PasswordGenerationConfig = {
    length: 12,
    uppercase: true,
    lowercase: true,
    numbers: true,
    symbols: true,
};

const RandomPasswordGeneration = (passwordGenerationConfig: PasswordGenerationConfig): string => {
    const lowercaseLetters = "abcdefghijklmnopqrstuvwxyz";
    const uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ";
    const numbers = "0123456789";
    const symbols = "!@#$%^&*()_+-=/~`[]{}|;:,.<>?";

    let password = "";
    for (let i = 0; i < passwordGenerationConfig.length; i++) {
        const activeOptions = [
            passwordGenerationConfig.uppercase && uppercaseLetters,
            passwordGenerationConfig.lowercase && lowercaseLetters,
            passwordGenerationConfig.numbers && numbers,
            passwordGenerationConfig.symbols && symbols,
        ].filter(opt => opt !== false);

        const randomOption = activeOptions[Math.floor(Math.random() * activeOptions.length)];
        
        const randomIndex = Math.floor(Math.random() * randomOption.length);
        password += randomOption[randomIndex];               
    }

    return password;
};

const copyToClipboard = (password: string) => {
    navigator.clipboard.writeText(password).then(() => {
        alert("Password copied to clipboard!");
    }).catch(err => {
        console.error('Could not copy text: ', err);
    });
};

const PasswordGeneratorComponent: React.FC = () => {

    const [password, setPassword] = React.useState<string>("");
    const [passwordGenerationConfig, setPasswordGenerationConfig] = React.useState<PasswordGenerationConfig>(DefaultConfig);

    const isAnyOptionChecked = React.useMemo(() => {

        return [passwordGenerationConfig.uppercase, 
         passwordGenerationConfig.lowercase, 
         passwordGenerationConfig.numbers, 
         passwordGenerationConfig.symbols].filter(Boolean).length > 1
    }, [passwordGenerationConfig]);

    const handleGeneratePassword = () => {
        setPassword(RandomPasswordGeneration(passwordGenerationConfig));
    };

    useEffect(() => {
        handleGeneratePassword();
    }, [passwordGenerationConfig]);

    const handleCheckboxChange = (key: keyof PasswordGenerationConfig) => (event: React.ChangeEvent<HTMLInputElement>) => {
        if (isAnyOptionChecked || event.target.checked === true) {
            setPasswordGenerationConfig({ ...passwordGenerationConfig, [key]: event.target.checked });
        } else {
            setPassword("Please select at least one option");
        }
    };

    const handleLengthChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setPasswordGenerationConfig({ ...passwordGenerationConfig, length: Number(event.target.value) });
    };

    return (
        <div className="flex flex-col items-center bg-white shadow-md rounded">
            <div className="flex w-full flex-col items-center p-3 rounded-xl shadow-lg">
                <h1 className="mb-4 text-3xl font-mono text-white bg-black px-4 py-2 rounded-md">
                    {password}
                    &nbsp;
                    <button 
                        type="button"
                        onClick={handleGeneratePassword}
                    >
                        🎲
                    </button>
                    &nbsp;
                    <button 
                        type="button"
                        onClick={() => copyToClipboard(password)}
                    >
                        ✂
                    </button>
                </h1>
                <div className="flex flex-row items-center w-full px-4 py-4 space-y-3 bg-gray-100 rounded-md shadow-inner">
                    <button 
                        className="bg-slate-800 hover:bg-slate-900 text-white font-bold text-xl py-2 w-full text-center rounded-r"
                        onClick={() => setPasswordGenerationConfig(DefaultConfig)}
                    >
                        Reset options
                    </button>
                    <div className="flex flex-col items-start w-full px-4 space-y-3 ">
                    
                        {['uppercase', 'lowercase', 'numbers', 'symbols'].map((type) => (
                            <label key={type} className="flex items-center space-x-2">
                                <input 
                                    type="checkbox" 
                                    className="h-5 w-5"
                                    checked={Boolean(passwordGenerationConfig[type as keyof PasswordGenerationConfig])} 
                                    onChange={handleCheckboxChange(type as keyof PasswordGenerationConfig)}
                                />
                                <span className="text-xl text-gray-800">{type.charAt(0).toUpperCase() + type.slice(1)}</span>
                            </label>
                        ))}
                    </div>
                    <div className="flex flex-col items-start w-full">
                        <div className="flex items-center justify-between space-x-4 w-full">
                            <span className="text-xl text-gray-800">Length:</span>
                            <input 
                                type="range" 
                                className="w-full accent-indigo-600"
                                min="1" 
                                max="50" 
                                value={passwordGenerationConfig.length} 
                                onChange={handleLengthChange}
                            />
                            <span className="text-lg text-gray-800">{passwordGenerationConfig.length}</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default PasswordGeneratorComponent;