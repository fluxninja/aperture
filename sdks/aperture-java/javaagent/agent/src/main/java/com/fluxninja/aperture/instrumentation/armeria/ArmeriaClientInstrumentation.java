package com.fluxninja.aperture.instrumentation.armeria;

import com.fluxninja.aperture.instrumentation.TransformerInstrumentation;
import net.bytebuddy.agent.builder.AgentBuilder;
import net.bytebuddy.description.type.TypeDescription;
import net.bytebuddy.matcher.ElementMatcher;
import net.bytebuddy.matcher.ElementMatchers;

import static net.bytebuddy.matcher.ElementMatchers.*;

public class ArmeriaClientInstrumentation implements TransformerInstrumentation {
    @Override
    public ElementMatcher<TypeDescription> getType() {
        return ElementMatchers.named("com.linecorp.armeria.client.WebClientBuilder");
    }

    @Override
    public AgentBuilder.Transformer getTransformer() {
        return new AgentBuilder.Transformer.ForAdvice()
                .advice(isMethod().and(isPublic()).and(named("build")),
                        ArmeriaClientAdvice.class.getName());
    }
}
