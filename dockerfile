FROM filvenus/venus-buildenv AS buildenv

COPY . ./venus-wallet
RUN export GOPROXY=https://goproxy.cn && cd venus-wallet  && make

RUN cd venus-wallet && ldd ./venus-wallet


FROM filvenus/venus-runtime

# DIR for app
WORKDIR /app

# copy the app from build env
COPY --from=buildenv  /go/venus-wallet/venus-wallet /app/venus-wallet



EXPOSE 5680

# run wallet app
ENTRYPOINT ["/app/venus-wallet","run"]
