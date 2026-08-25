export const useAuth = () => {
  const login = async (credentials: {
    email: string
    password: string
  }) => {
    console.log('Login:', credentials)

    // Later:
    // await $fetch('/auth/login', {
    //   method: 'POST',
    //   body: credentials
    // })
  }

  const register = async (credentials: {
    name: string
    email: string
    password: string
  }) => {
    console.log('Register:', credentials)

    // Later:
    // await $fetch('/auth/register', {
    //   method: 'POST',
    //   body: credentials
    // })
  }

  return {
    login,
    register
  }
}