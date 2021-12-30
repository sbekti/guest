import {
  Box,
  FormControl,
  FormLabel,
  Input,
  Checkbox,
  Stack,
  Button,
  Link,
  Image,
  FormErrorMessage,
  Alert,
  AlertIcon,
  useColorModeValue,
} from '@chakra-ui/react';
import { useEffect, useState } from 'react';
import { Link as RouteLink } from 'react-router-dom';
import { Formik, Form, Field, ErrorMessage } from 'formik';
import * as Yup from 'yup';

export const RegisterCard = ({ onRegisterSuccess }) => {
  const [error, setError] = useState(null);
  const [captchaId, setCaptchaId] = useState('');
  const [captchaUrl, setCaptchaUrl] = useState('');

  const loadCaptcha = () => {
    fetch('/api/v1/captcha')
      .then(res => res.json())
      .then(
        result => {
          setCaptchaId(result.captcha_id);
          setCaptchaUrl(`/captcha/${result.captcha_id}.png`);
        },
        error => {
          setError(error);
        }
      );
  };

  const handleSubmit = (
    values,
    setSubmitting,
    resetForm,
    setErrors,
    setFieldValue
  ) => {
    fetch('/api/v1/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        email: values.email,
        captcha_id: captchaId,
        captcha_answer: values.captchaAnswer,
        corp_access: values.corpAccess,
      }),
    })
      .then(res => res.json())
      .then(
        result => {
          if (result.success) {
            onRegisterSuccess(values.corpAccess);
            resetForm();
          } else {
            setError(result.message);
            setErrors(result.input_errors);
            if (result.input_errors['captchaAnswer']) {
              setFieldValue('captchaAnswer', '', false);
            }
            loadCaptcha();
          }
          setSubmitting(false);
        },
        error => {
          setError(error);
          loadCaptcha();
          setSubmitting(false);
        }
      );
  };

  useEffect(() => {
    loadCaptcha();
  }, []);

  return (
    <Box
      rounded={'lg'}
      bg={useColorModeValue('white', 'gray.700')}
      borderWidth="1px"
      w={80}
      p={4}
    >
      <Stack spacing={6}>
        {error && (
          <Alert status="error">
            <AlertIcon />
            {error}
          </Alert>
        )}
        <Formik
          initialValues={{
            email: '',
            captchaAnswer: '',
            corpAccess: false,
            agreed: false,
          }}
          validationSchema={Yup.object({
            email: Yup.string()
              .email('Invalid email address')
              .required('Required'),
            captchaAnswer: Yup.string().required('Required'),
            agreed: Yup.boolean().oneOf([true], 'Must be checked'),
          })}
          onSubmit={(
            values,
            { setSubmitting, resetForm, setErrors, setFieldValue }
          ) => {
            handleSubmit(
              values,
              setSubmitting,
              resetForm,
              setErrors,
              setFieldValue
            );
          }}
        >
          {({ values, isSubmitting, isValid, dirty, setFieldValue }) => (
            <Form>
              <Stack spacing={4}>
                <Field type="email" name="email">
                  {({ field, meta }) => (
                    <FormControl isInvalid={meta.touched && meta.error}>
                      <FormLabel>Email address</FormLabel>
                      <Input type="email" {...field} />
                      <ErrorMessage name="email" component={FormErrorMessage} />
                    </FormControl>
                  )}
                </Field>
                <Field type="text" name="captchaAnswer">
                  {({ field, meta }) => (
                    <FormControl isInvalid={meta.touched && meta.error}>
                      <FormLabel w="100%">
                        Type the numbers you see
                        <Box align="center" mt={1} mb={1}>
                          <Image
                            h={20}
                            w={60}
                            src={captchaUrl}
                            alt="CAPTCHA image"
                          />
                        </Box>
                      </FormLabel>
                      <Input type="text" inputMode="numeric" {...field} />
                      <Input
                        type="hidden"
                        name="captchaId"
                        defaultValue={captchaId}
                      />
                      <ErrorMessage
                        name="captchaAnswer"
                        component={FormErrorMessage}
                      />
                    </FormControl>
                  )}
                </Field>
                <Stack spacing={10}>
                  <Stack>
                    <Field type="checkbox" name="corpAccess">
                      {() => (
                        <Checkbox
                          isChecked={values.corpAccess}
                          onChange={e =>
                            setFieldValue('corpAccess', e.target.checked)
                          }
                        >
                          Request corp network access
                        </Checkbox>
                      )}
                    </Field>
                    <Field type="checkbox" name="agreed">
                      {() => (
                        <Checkbox
                          isChecked={values.agreed}
                          onChange={e =>
                            setFieldValue('agreed', e.target.checked)
                          }
                        >
                          I agree to the{' '}
                          <Link color={'blue.400'} as={RouteLink} to="/terms">
                            Terms of Use
                          </Link>
                        </Checkbox>
                      )}
                    </Field>
                  </Stack>
                  <Button
                    type="submit"
                    isDisabled={isSubmitting || !(isValid && dirty)}
                    isLoading={isSubmitting}
                    bg={'blue.400'}
                    color={'white'}
                    _hover={{
                      bg: 'blue.500',
                    }}
                  >
                    Register
                  </Button>
                </Stack>
              </Stack>
            </Form>
          )}
        </Formik>
      </Stack>
    </Box>
  );
};
