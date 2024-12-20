import PasswordCheckerComponent from './component/password-checker/password-checker';
import PasswordGeneratorComponent from './component/password-generator/password-generator';
import WebService from './service/web/web-service';

function App() {

  const webService = new WebService(import.meta.env.VITE_BACKEND_URL);

  return (
    <>
      <PasswordCheckerComponent webService={webService} />
      <PasswordGeneratorComponent />
    </>
  )
}

export default App
