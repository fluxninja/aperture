package com.fluxninja.aperture.instrumentation.netty;

import com.fluxninja.aperture.instrumentation.TransformerInstrumentation;
import net.bytebuddy.agent.builder.AgentBuilder;
import net.bytebuddy.description.type.TypeDescription;
import net.bytebuddy.matcher.ElementMatcher;
import net.bytebuddy.matcher.ElementMatchers;

import static net.bytebuddy.matcher.ElementMatchers.*;

public class NettyServerInstrumentation implements TransformerInstrumentation {
    @Override
    public ElementMatcher<TypeDescription> getType() {
        return ElementMatchers.named("io.netty.channel.DefaultChannelPipeline");
    }

    @Override
    public AgentBuilder.Transformer getTransformer() {
        return new AgentBuilder.Transformer.ForAdvice()
                .advice(isMethod().and(nameStartsWith("add"))
                                .and(takesArgument(1, String.class))
                                .and(takesArgument(2, named("io.netty.channel.ChannelHandler"))),
                        NettyServerAdvice.class.getName());
    }

}
