export const mode = process.env.NODE_ENV || 'development';

export const apiHost = mode === 'production' ? 'https://api.linkinpark.tylerseabury.com' : 'http://localhost:7777';